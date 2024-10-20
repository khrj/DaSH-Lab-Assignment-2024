package main

import (
	"context"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	ctx    context.Context
	client *genai.Client
	model  *genai.GenerativeModel
)

func init_ai() {
	ctx = context.Background()

	apiKey, ok := os.LookupEnv("GEMINI_API_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	var err error
	client, err = genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	model = client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(1)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"
}

func ask(q string) string {
	log.Println("⏳ Processing: ", q)

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(q))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	var responseText string

	for _, part := range resp.Candidates[0].Content.Parts {
		if s, ok := part.(genai.Text); ok {
			responseText += string(s)
		}
	}

	log.Println("✅ Processed: ", q)

	return responseText
}
