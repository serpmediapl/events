package ocr

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/otiai10/gosseract"
	"github.com/serpmediapl/home-finance/pkg/events"
	bolt "go.etcd.io/bbolt"
)

// EventTypeOcr states ocr event type
const EventTypeOcr = events.EventType("ocr")

type uploadData struct {
	FilePath string `json:"file_path"`
}

type ocrData struct {
	Text string `json:"text"`
}

type EventHandler struct {
	dbPath string
	d      events.Dispatcher
}

func NewEventHandler(dbPath string, d events.Dispatcher) *EventHandler {
	return &EventHandler{dbPath: dbPath, d: d}
}

func (h *EventHandler) Handle(e events.Event) error {
	data := &uploadData{}
	err := json.Unmarshal(e.Data, data)
	if err != nil {
		return err
	}
	img, err := os.Open(data.FilePath)
	if err != nil {
		return err
	}
	client := gosseract.NewClient()
	defer client.Close()
	var buf bytes.Buffer
	_, err = buf.ReadFrom(img)
	if err != nil {
		return err
	}
	err = client.SetImageFromBytes(buf.Bytes())
	if err != nil {
		return err
	}
	client.SetLanguage("pol")
	out, err := client.Text()
	if err != nil {
		return err
	}
	chnks := strings.Split(data.FilePath, "/")
	err = h.d.Send(events.Event{
		TimeStamp: time.Now().UTC(),
		Type:      EventTypeOcr,
		User:      chnks[len(chnks)-2],
		Data:      json.RawMessage(out),
	})
	if err != nil {
		return err
	}
	return saveText(h.dbPath, out, data.FilePath)
}

func saveText(dbName, text string, key string) error {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("texts"))
		return b.Put([]byte(key), []byte(text))
	})
}
