// internal/infrastructure/event/event_bus.go
package event

type Event[T any] struct {
	Payload T
}

type EventHandler[T any] func(Event[T])

type EventBus[T any] struct {
	handlers []EventHandler[T]
}

func NewEventBus[T any]() *EventBus[T] {
	return &EventBus[T]{handlers: []EventHandler[T]{}}
}

func (b *EventBus[T]) Publish(e Event[T]) {
	for _, h := range b.handlers {
		go h(e) // 异步调用
	}
}

func (b *EventBus[T]) Subscribe(h EventHandler[T]) {
	b.handlers = append(b.handlers, h)
}
