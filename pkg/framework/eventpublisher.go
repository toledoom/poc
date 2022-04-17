package framework

import "github.com/hashicorp/go-multierror"

type Event interface {
	Name() string
}

type EventHandler interface {
	Notify(event Event) error
}

type EventPublisher struct {
	handlers map[string][]EventHandler
}

func NewEventPublisher() *EventPublisher {
	return &EventPublisher{
		handlers: make(map[string][]EventHandler),
	}
}

func (e *EventPublisher) Subscribe(handler EventHandler, events ...Event) {
	for _, event := range events {
		handlers := e.handlers[event.Name()]
		handlers = append(handlers, handler)
		e.handlers[event.Name()] = handlers
	}
}

func (e *EventPublisher) Notify(event Event) error {
	var multipleError error
	n := event.Name()
	for _, handler := range e.handlers[n] {
		err := handler.Notify(event)
		if err != nil {
			multipleError = multierror.Append(multipleError, err)
		}
	}
	return multipleError
}
