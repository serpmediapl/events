package events

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

// NatsDispatcher dispatches events through nats queues
type NatsDispatcher struct {
	URL string
}

// NewNatsDispatcher creates NatsDispatcher instance
func NewNatsDispatcher(URL string) *NatsDispatcher {
	return &NatsDispatcher{URL: URL}
}

// Send implements event dispatcher interface
func (nd *NatsDispatcher) Send(e Event) error {
	nc, err := natsCli(nd.URL)
	if err != nil {
		return err
	}
	defer nc.Close()
	return nc.Publish("events", e)
}

// NatsSink receives events through nats queue
type NatsSink struct {
	URL   string
	errCh chan error
}

// NewNatsSink creates NatsSink instance
func NewNatsSink(natsURL string) *NatsSink {
	return &NatsSink{URL: natsURL}
}

// Receive receives events
func (ns *NatsSink) Receive() chan Event {
	out := make(chan Event)
	go func(out chan Event, natsURL string, ns *NatsSink) {
		defer close(out)
		defer close(ns.errCh)
		nc, err := natsCli(natsURL)
		if err != nil {
			ns.errCh <- err
			return
		}
		defer nc.Close()
		ch := make(chan *nats.Msg)
		sub, err := nc.Conn.ChanSubscribe("events", ch)
		if err != nil {
			ns.errCh <- err
			return
		}
		defer func(sub *nats.Subscription, ch chan *nats.Msg) {
			sub.Unsubscribe()
			close(ch)
		}(sub, ch)
		for msg := range ch {
			var e Event
			err := json.Unmarshal(msg.Data, &e)
			if err != nil {
				ns.errCh <- err
				continue
			}
			out <- e
		}
	}(out, ns.URL, ns)
	return out
}

// Err returns any error that occured during processing
func (ns *NatsSink) Err() chan error {
	return ns.errCh
}

func natsCli(natsURL string) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return nats.NewEncodedConn(nc, nats.JSON_ENCODER)
}
