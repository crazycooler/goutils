package eventbus

import (
	"log"
	"sync"

	"github.com/crazycooler/goutils/utils"
)

type EventItem struct {
	Payload string
	Topic   string
}

type EventChannel chan *EventItem

// 事件总线
type EventBus interface {
	Subscribe(topic string, ch EventChannel)   //订阅消息
	UnSubscribe(topic string, ch EventChannel) //取消订阅
	Publish(topic string, payload string)      //发布消息

	Trigger(topic string, item *EventItem) //用于触发消息订阅事件
	Close()                                //关闭
}

// 发布接口
type EventPublisher interface {
	Publish(topic string, item *EventItem, bus EventBus) error
	Close()
}

type EventBusBase struct {
	subscriber map[string][]EventChannel //订阅者
	lock       sync.RWMutex
	publisher  EventPublisher //发布者
}

// 订阅消息
func (s *EventBusBase) Subscribe(topic string, ch EventChannel) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if slice, ok := s.subscriber[topic]; ok {
		index := utils.ArrayIndex(ch, slice)
		if index != -1 {
			//如果已经订阅了，则不重复订阅
			return
		}
		s.subscriber[topic] = append(slice, ch)
	} else {
		s.subscriber[topic] = []EventChannel{ch}
	}

}

// 取消订阅
func (s *EventBusBase) UnSubscribe(topic string, ch EventChannel) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if slice, ok := s.subscriber[topic]; ok {
		index := utils.ArrayIndex(ch, slice)
		if index == -1 {
			return
		}

		s.subscriber[topic] = append(slice[:index], slice[index+1:]...)
	}
}

// 发布消息
func (s *EventBusBase) Publish(topic string, payload string) {
	if s.publisher == nil {
		return
	}
	item := &EventItem{
		Topic:   topic,
		Payload: payload,
	}
	s.publisher.Publish(topic, item, s)
}

// 用于触发消息订阅事件
// func (s *EventBusBase) Trigger(topic string, item *EventItem) {
// 	s.lock.RLock()
// 	defer s.lock.RUnlock()

// 	if slice, ok := s.subscriber[topic]; ok {
// 		//新建的newSlice不会被, 不会被上面修改
// 		newSlice := append([]EventChannel{}, slice...)
// 		go func(item *EventItem, newSlice []EventChannel) {
// 			for _, ch := range newSlice {
// 				select {
// 				case ch <- item:
// 				default:
// 					//如果通道中消息已满，则丢弃，防止阻塞，导致协程泄露。
// 					log.Println("write chan full!!!")
// 				}
// 			}
// 		}(item, newSlice)
// 	}
// }

// 用于触发消息订阅事件, 采用同步方式发送
func (s *EventBusBase) Trigger(topic string, item *EventItem) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if slice, ok := s.subscriber[topic]; ok {
		//新建的newSlice不会被, 不会被上面修改
		for _, ch := range slice {
			select {
			case ch <- item:
			default:
				//如果通道中消息已满，则丢弃，防止阻塞。
				log.Println("write chan full!!!")
			}
		}
	}
}

// 关闭
func (s *EventBusBase) Close() {
	if s.publisher != nil {
		s.publisher.Close()
		s.publisher = nil
	}
}
