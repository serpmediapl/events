package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/serpmediapl/home-finance/pkg/events"
)

// EventTypeUpload points to event type that this pkg will publish
const EventTypeUpload = "upload"

// DefaultRoot points to default image storage directory
const DefaultRoot = "/images"

// DefaultLogger logs to stderr with default flags
var DefaultLogger = log.New(os.Stderr, "", log.LstdFlags)

type uploadData struct {
	FilePath string `json:"file_path"`
}

// Handler handles upload receipt image
type Handler struct {
	dp   events.Dispatcher
	log  *log.Logger
	root string
}

// NewHandler creates new image upload handler
func NewHandler() *Handler {
	return &Handler{root: DefaultRoot, log: DefaultLogger}
}

// WithRoot sets image root for this handler
func (h *Handler) WithRoot(root string) *Handler {
	h.root = root
	return h
}

// WithEventDispatcher sets event dispatcher for this handler
func (h *Handler) WithEventDispatcher(dp events.Dispatcher) *Handler {
	h.dp = dp
	return h
}

// WithLogger sets logger for this handler
func (h *Handler) WithLogger(l *log.Logger) *Handler {
	h.log = l
	return h
}

// ServeHTTP implements http.Handler interface
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var buf bytes.Buffer
	tee := io.TeeReader(r.Body, &buf)
	hs := sha256.New()
	_, err := io.Copy(hs, tee)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	username := r.Header.Get("X-WebAuth-User")
	fileName := fmt.Sprintf("%x.jpg", hs.Sum(nil))
	filePath := filepath.Join(h.root, username, fileName)
	_, err = os.Stat(filePath)
	if err == nil {
		http.Error(w, http.StatusText(409), 409)
		return
	}
	img, err := os.Create(filePath)
	if err != nil {
		h.log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	h.log.Printf("Saving file: %s", filePath)
	_, err = io.Copy(img, &buf)
	if err != nil {
		h.log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	var evBuf bytes.Buffer
	err = json.NewEncoder(&evBuf).Encode(&uploadData{FilePath: filePath})
	if err != nil {
		h.log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	event := events.Event{
		Type:      EventTypeUpload,
		User:      username,
		TimeStamp: time.Now().UTC(),
		Data:      json.RawMessage(evBuf.Bytes()),
	}
	err = h.dp.Send(event)
	if err != nil {
		h.log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(204)
}
