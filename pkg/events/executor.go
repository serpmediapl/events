package events

import "errors"

var ErrNotRegistered = errors.New("Event handler not registered")

type EventType string

type Executor interface {
	Execute(Event) error
	Register(EventType, Handler) error
}

type InMemoryExecutor struct {
	db map[EventType]Handler
}

func NewInMemoryExecutor() *InMemoryExecutor {
	return &InMemoryExecutor{db: make(map[EventType]Handler)}
}

func (r *InMemoryExecutor) Register(et EventType, eh Handler) error {
	r.db[et] = eh
	return nil
}

func (r *InMemoryExecutor) Execute(e Event) error {
	handler, ok := r.db[e.Type]
	if !ok {
		return ErrNotRegistered
	}
	return handler.Handle(e)
}
