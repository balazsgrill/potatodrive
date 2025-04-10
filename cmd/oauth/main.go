package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig = &oauth2.Config{
		Scopes:   []string{"https://www.googleapis.com/auth/photoslibrary"},
		Endpoint: google.Endpoint,
	}
	state = uuid.New().String()
)

func main() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: main <clientID> <clientSecret> <redirectURL>")
	}

	oauthConfig.ClientID = os.Args[1]
	oauthConfig.ClientSecret = os.Args[2]
	oauthConfig.RedirectURL = os.Args[3]

	fmt.Printf("Starting server at %s\n", oauthConfig.RedirectURL)
	u, err := url.Parse(oauthConfig.RedirectURL)
	if err != nil {
		log.Fatalf("Invalid redirect URL: %v", err)
	}
	port := u.Port()
	if port == "" {
		if u.Scheme == "https" {
			port = "443"
		} else if u.Scheme == "http" {
			port = "80"
		} else {
			log.Fatalf("Unknown scheme: %s", u.Scheme)
		}
		log.Fatal("No port specified in redirect URL")
	}

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	callbackPath := "/callback"
	if u.Path != "" && u.Path != "/" {
		callbackPath = u.Path
	}
	http.HandleFunc(callbackPath, handleCallback)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<html><body><a href="/login">Login with Google</a></body></html>`)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := os.Create("token.json")
	if err != nil {
		http.Error(w, "Failed to create token file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(token); err != nil {
		http.Error(w, "Failed to save token to file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(fmt.Sprintf("User Info: %s", resp.Body))); err != nil {
		log.Println("Failed to write response:", err)
	}
}
