package eventbus

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// 初始化redis客户端
func InitializeRedis(redisUrl string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redisUrl,
		Password: "", // no password set
		DB:       1,  // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Panicln("Redis connect ping failed, err:", zap.Any("err", err))
		return nil
	}

	return client
}

func NewRedisEventBus(client *redis.Client, redisTopic string) EventBus {
	if client == nil {
		return nil
	}

	pubsub := client.Subscribe(context.Background(), redisTopic)
	ctx, cancel := context.WithCancel(context.Background())

	publisher := &EventRedisPublisher{
		client:     client,
		pubsub:     pubsub,
		cancel:     cancel,
		redisTopic: redisTopic,
	}

	bus := &EventBusBase{
		subscriber: make(map[string][]EventChannel),
		publisher:  publisher,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("RedisEventBus done !!!")
				return
			default:
				msg, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					log.Printf("ReceiveMessage error: %v", err)
					continue
				}

				var item EventItem
				err = json.Unmarshal([]byte(msg.Payload), &item)
				if err != nil {
					log.Printf("bad redis payload: %s\n", msg.Payload)
					continue
				}

				//获取到redis的消息后，用来触发Event bus
				bus.Trigger(item.Topic, &item)

			}
		}
	}()

	return bus
}

type EventRedisPublisher struct {
	client     *redis.Client //redis client
	pubsub     *redis.PubSub //redis pubsub
	cancel     context.CancelFunc
	redisTopic string //在redis中订阅的topic，和EventBus中的topic不是同一个
}

func (s *EventRedisPublisher) Publish(topic string, item *EventItem, bus EventBus) error {
	buf, _ := json.Marshal(item) //用json做序列化
	//推送消息给redis
	err := s.client.Publish(context.Background(), s.redisTopic, buf).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *EventRedisPublisher) Close() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	if s.pubsub != nil {
		s.pubsub.Close()
		s.pubsub = nil
	}
}
