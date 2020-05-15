package events

import (
	"encoding/json"
	"time"
)

// Event represents application event
type Event struct {
	TimeStamp time.Time       `json:"time_stamp"`
	Type      EventType       `json:"type"`
	Data      json.RawMessage `json:"message"`
	User      string          `json:"user"`
}

// Dispatcher sends event to subscribers
type Dispatcher interface {
	Send(Event) error
}

// Sink receives event stream
type Sink interface {
	Receive() chan Event
	Err() chan error
}
