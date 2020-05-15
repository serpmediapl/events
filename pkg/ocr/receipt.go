package ocr

import (
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	ID        uuid.UUID `json:"id"`
	Items     []Item    `json:"items"`
	Store     string    `json:"store"`
	DateTime  time.Time `json:"date_time"`
	Timestamp time.Time `json:"timestamp"`
}

type Item struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}
