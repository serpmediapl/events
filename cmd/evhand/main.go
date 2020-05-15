package main

import (
	"flag"
	"log"
	"os"

	"github.com/serpmediapl/home-finance/pkg/events"
	"github.com/serpmediapl/home-finance/pkg/ocr"
	"github.com/serpmediapl/home-finance/pkg/parser"
)

func main() {
	logger := log.New(os.Stdout, "[evhand] ", log.LstdFlags|log.Lshortfile)
	ocrDbPath := flag.String("ocrdbpath", "ocr.db", "OCR Database path")
	natsURL := flag.String("natsurl", "", "Nats url")
	flag.Parse()

	dispatcher := events.NewNatsDispatcher(*natsURL)

	exec := events.NewInMemoryExecutor()
	exec.Register(ocr.EventTypeOcr, ocr.NewEventHandler(*ocrDbPath, dispatcher))
	exec.Register(parser.EventTypeParse, parser.NewEventHandler(dispatcher))

	sink := events.NewNatsSink(*natsURL)
	recCh := sink.Receive()
	errCh := sink.Err()

	for {
		select {
		case e := <-recCh:
			logger.Printf("Executing event: %s", e.Type)
			err := exec.Execute(e)
			if err != nil {
				logger.Println(err)
			}
		case err := <-errCh:
			logger.Printf("Error executing event: %+v", err)
		}

		if recCh == nil && errCh == nil {
			break
		}
	}
}
