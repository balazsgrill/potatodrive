package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getClient(config *oauth2.Config) *http.Client {
	// Load token from file
	tokenFile := "token.json"
	file, err := os.Open(tokenFile)
	if err != nil {
		log.Fatalf("Unable to open token file: %v", err)
	}
	defer file.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)
	if err != nil {
		log.Fatalf("Unable to decode token: %v", err)
	}

	// Create HTTP client
	return config.Client(context.Background(), token)
}

func main() {
	// Load OAuth2 configuration
	config, err := google.ConfigFromJSON([]byte(`YOUR_CREDENTIALS_JSON`), "https://www.googleapis.com/auth/drive")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	// Use the client (example)
	resp, err := client.Get("https://www.googleapis.com/drive/v3/files")
	if err != nil {
		log.Fatalf("Unable to make request: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Request successful")
}
