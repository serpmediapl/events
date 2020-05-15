package ocr

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

var stores = []string{"stokrotka", "lidl", "biedronka"}

func Parse(ctx context.Context, r io.Reader) Receipt {
	reItem, _ := regexp.Compile(`[\d,]*\ \*\ [\d,]*`)
	scanner := bufio.NewScanner(r)
	rec := Receipt{
		ID:        uuid.New(),
		Items:     []Item{},
		Timestamp: time.Now().UTC(),
	}
	for scanner.Scan() {
		line := scanner.Bytes()
		if rec.Store == "" {
			for _, store := range stores {
				if parseStore(line, store) {
					rec.Store = store
					continue
				}
			}
		}
		if rec.DateTime.IsZero() {
			rec.DateTime, _ = time.Parse("2006-01-02", string(line))
			continue
		}
		idxs := reItem.FindIndex(line)
		if idxs != nil {
			rec.Items = append(rec.Items, parseItem(line, idxs))
		}
	}
	return rec
}

func parseItem(line []byte, idxs []int) Item {
	name := line[0 : idxs[0]-1]
	amnt := line[idxs[0]:idxs[1]]
	amnt = bytes.ReplaceAll(amnt, []byte(","), []byte("."))
	chnks := bytes.Split(amnt, []byte(" * "))
	quantity, err := strconv.ParseFloat(string(chnks[0]), 32)
	if err != nil {
		log.Println(err)
	}
	amount, err := strconv.ParseFloat(string(chnks[1]), 32)
	if err != nil {
		log.Println(err)
	}
	return Item{
		Name:     string(name),
		Price:    math.Round(amount*100) / 100,
		Quantity: math.Round(quantity*100) / 100,
	}
}

func parseStore(line []byte, store string) bool {
	chnks := strings.Split(string(line), " ")
	for _, chnk := range chnks {
		if strings.ToLower(chnk) == store {
			return true
		}
	}
	return false
}
