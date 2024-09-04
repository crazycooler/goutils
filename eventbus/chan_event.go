package eventbus

func NewChanEventBus() EventBus {
	publisher := &EventChanPublisher{}
	bus := &EventBusBase{
		subscriber: make(map[string][]EventChannel),
		publisher:  publisher,
	}

	return bus
}

type EventChanPublisher struct {
}

func (s *EventChanPublisher) Publish(topic string, item *EventItem, bus EventBus) error {
	bus.Trigger(topic, item)
	return nil
}

func (s *EventChanPublisher) Close() {

}
