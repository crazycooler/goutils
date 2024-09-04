package eventbus

import (
	"testing"
	"time"
)

func TestRedisEventBus(t *testing.T) {
	client := InitializeRedis("127.0.0.1:6379")
	bus := NewRedisEventBus(client, "redis_topic_0904")

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
