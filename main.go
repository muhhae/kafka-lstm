package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Data struct {
	Value     float32   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	addr := os.Getenv("ADDR")
	writers := []*kafka.Writer{
		{
			Addr:                   kafka.TCP(addr),
			Topic:                  "temp_1",
			AllowAutoTopicCreation: true,
		}, {
			Addr:                   kafka.TCP(addr),
			Topic:                  "humid_1",
			AllowAutoTopicCreation: true,
		}, {
			Addr:                   kafka.TCP(addr),
			Topic:                  "light_intensity",
			AllowAutoTopicCreation: true,
		},
	}

	url := "https://ssiot.jlbsd.my.id/api/all"
	timeLayout := "2006-01-02 15:04:05"

	now := time.Now().UnixMilli()
	var prev int64 = 0

	for {
		if now = time.Now().UnixMilli(); now-prev < 1000 {
			continue
		}
		prev = now
		response, err := http.Get(url)

		if response.StatusCode != http.StatusOK {
			log.Printf("Failed To Fetch Error: ", response.StatusCode)
			continue
		}

		defer response.Body.Close()

		var result []Data
		err = json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			log.Printf("Failed to Decode Error: ", err.Error())
			continue
		}
		fmt.Println(result)

		for i, w := range writers {
			err = w.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte("value"),
					Value: []byte(fmt.Sprintf("%f", result[i*2].Value)),
				},
				kafka.Message{
					Key:   []byte("timestamp"),
					Value: []byte(result[i].Timestamp.Format(timeLayout)),
				},
			)
			if err != nil {
				log.Printf("Failed to Publish: ", err.Error())
				break
			}
		}
	}
}
