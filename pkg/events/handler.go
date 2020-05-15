package events

// Handle describes event handler behavior
type Handler interface {
	Handle(Event) error
}
