package main

import (
	"context"
	"log"
	"os"

	"github.com/balazsgrill/potatodrive/gpfs"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig = &oauth2.Config{
		Scopes:   []string{"https://www.googleapis.com/auth/photoslibrary"},
		Endpoint: google.Endpoint,
	}
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: main <clientID> <clientSecret> <redirectURL>")
	}

	oauthConfig.ClientID = os.Args[1]
	oauthConfig.ClientSecret = os.Args[2]
	oauthConfig.RedirectURL = os.Args[3]

	_, err := gpfs.Authenticate(context.Background(), oauthConfig.ClientID, oauthConfig.ClientSecret, oauthConfig.RedirectURL, func(s string) {
		log.Println("Open this URL in your browser:", s)
	})
	if err != nil {
		log.Fatalf("Error during authentication: %v", err)
	}
}
