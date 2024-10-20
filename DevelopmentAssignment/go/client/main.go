package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
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
	ClientId     uint64
	Source       string
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalln("Error connecting to NATS:", err)
	}

	inputData, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	questions := strings.Split(string(inputData), "\n")

	outputs := []Output{}

	nc.Subscribe("response", func(m *nats.Msg) {
		resp := messageToResponse(m)

		source := "user"
		if slices.Contains(questions, resp.Prompt) {
			source = "gemini"
		}

		clientId, err := nc.GetClientID()
		if err != nil {
			clientId = 0
		}

		output := Output{
			Prompt:       resp.Prompt,
			Message:      resp.Message,
			TimeSent:     resp.TimeSent,
			TimeReceived: resp.TimeReceived,
			ClientId:     clientId,
			Source:       source,
		}

		outputs = append(outputs, output)

		if len(outputs) >= 12 {
			byte_outputs, err := json.Marshal(outputs)

			if err != nil {
				log.Println("Error marshalling outputs: ", err)
				return
			}

			err = os.WriteFile("outputs/output-"+strconv.FormatUint(clientId, 10)+".json", byte_outputs, 0644)

			if err != nil {
				log.Fatalln("Error writing output file: ", err)
			}

			os.Exit(0)
		}
	})

	for _, question := range questions {
		time.Sleep(4 * time.Second)

		query := Query{
			Prompt:   question,
			TimeSent: time.Now(),
		}

		log.Println("Sending query:", query)
		nc.Publish("query", toGob(query))
	}

	select {} // block forever
}

func messageToResponse(m *nats.Msg) Response {
	var response Response

	err := gob.NewDecoder(bytes.NewReader(m.Data)).Decode(&response)
	if err != nil {
		log.Fatalln("Error decoding query:", err)
	}

	return response
}

func toGob(value any) []byte {
	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(value)

	if err != nil {
		log.Fatalln("Error encoding query:", err)
	}

	return buf.Bytes()
}
