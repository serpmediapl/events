package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/serpmediapl/home-finance/pkg/events"
)

func main() {
	logger := log.New(os.Stdout, "[recupl] ", log.LstdFlags|log.Lshortfile)
	natsURL := flag.String("natsurl", "", "Nats url")
	port := flag.Int("port", 8000, "Port to listen to")
	flag.Parse()

	h := NewHandler().
		WithEventDispatcher(events.NewNatsDispatcher(*natsURL)).
		WithLogger(logger)

	addr := fmt.Sprintf(":%d", *port)
	logger.Fatal(http.ListenAndServe(addr, h))
}
