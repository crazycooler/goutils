package eventbus

import (
	"testing"
	"time"
)

func TestChanEventBus(t *testing.T) {
	bus := NewChanEventBus()

	ch := make(EventChannel, 64)
	bus.Subscribe("topic", ch)

	bus.Publish("topic", "hello world!!!")
	bus.Publish("topic", "hello golang!!!")

	go func() {
		select {
		case item := <-ch:
			if item.Topic != "topic" || item.Payload != "hello world!!!" {
				t.Errorf("expected, topic: topic, payload: hello world!!!, but got: %#v", item)
			}
		default:
			t.Errorf("expected get event item, but null")
		}

		select {
		case item := <-ch:
			if item.Topic != "topic" || item.Payload != "hello golang!!!" {
				t.Errorf("expected, topic: topic, payload: hello golang!!!, but got: %#v", item)
			}
		default:
			t.Errorf("expected get event item, but null")
		}
	}()

	time.Sleep(time.Second)
}
