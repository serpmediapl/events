package parser

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/serpmediapl/home-finance/pkg/events"
	"github.com/serpmediapl/home-finance/pkg/ocr"
)

const EventTypeParse = events.EventType("parse")

type ocrData struct {
	Text string `json:"text"`
}

type EventHandler struct {
	d events.Dispatcher
}

func NewEventHandler(d events.Dispatcher) *EventHandler {
	return &EventHandler{d: d}
}

func (eh *EventHandler) Handle(e events.Event) error {
	data := &ocrData{}
	err := json.Unmarshal(e.Data, data)
	if err != nil {
		return err
	}
	buf := strings.NewReader(data.Text)
	ocr.Parse(context.Background(), buf)
	return nil
}
