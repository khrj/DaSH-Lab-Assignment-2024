package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nats-io/nats.go"
)

type Query struct {
	Prompt   string
	TimeSent time.Time
}

type Response struct {
	Prompt       string
	TimeSent     time.Time
	Message      string
	TimeReceived time.Time
}

type Output struct {
	Prompt       string
	Message      string
	TimeSent     time.Time
	TimeReceived time.Time
	ClientId     string
	Source       string
}

func main() {
	init_ai()

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln("Error connecting to NATS:", err)
	}

	nc.Subscribe("query", func(m *nats.Msg) {
		go func() {
			query := messageToQuery(m)

			resp := Response{
				Prompt:       query.Prompt,
				TimeSent:     query.TimeSent,
				Message:      ask(query.Prompt),
				TimeReceived: time.Now(),
			}

			nc.Publish("response", toGob(resp))
		}()
	})

	select {} // block forever
}

func messageToQuery(m *nats.Msg) Query {
	var query Query

	err := gob.NewDecoder(bytes.NewReader(m.Data)).Decode(&query)
	if err != nil {
		log.Fatalln("Error decoding query:", err)
	}

	return query
}

func toGob(value any) []byte {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(value)

	if err != nil {
		log.Fatalln("Error encoding query:", err)
	}

	return buf.Bytes()
}
