package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/balazsgrill/potatodrive/bindings/proxy/server"
	"github.com/spf13/afero"
)

var Version string = "0.0.0-dev"

func main() {
	addr := flag.String("addr", ":8080", "The address to listen on for HTTP requests")
	directory := flag.String("directory", ".", "The directory to serve files from")
	keyID := flag.String("key-id", "", "The key ID for authentication")
	secret := flag.String("secret", "", "The secret for authentication")
	help := flag.Bool("help", false, "Show help message")
	ver := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *ver {
		fmt.Printf("Version: %s\n", Version)
		os.Exit(0)
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), *directory)
	mux := http.NewServeMux()

	// Wrap the handler with authentication middleware
	authenticatedHandler := authMiddleware(*keyID, *secret, server.Handler(fs))
	mux.HandleFunc("/", authenticatedHandler)

	httpserver := http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	httpserver.ListenAndServe()
}

func authMiddleware(keyID, secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providedKeyID := r.Header.Get("X-Key-ID")
		providedSecret := r.Header.Get("X-Secret")

		if providedKeyID != keyID || providedSecret != secret {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
