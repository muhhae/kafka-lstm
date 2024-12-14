package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	addr := "localhost:9092"
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{addr},
		Topic:     "humid_1",
		Partition: 0,
		MaxBytes:  10e6,
	})
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
