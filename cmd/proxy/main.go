package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/balazsgrill/potatodrive/bindings/proxy/server"
)

var Version string = "0.0.0-dev"

func main() {
	addr := flag.String("addr", ":8080", "The address to listen on for HTTP requests")
	configfile := flag.String("config", "", "Path to the configuration file")
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

	if *configfile == "" {
		fmt.Println("Please provide a configuration file")
		os.Exit(1)
	}

	log.Info().Str("version", Version).Msg("Starting PotatoDrive Proxy")
	config := loadConfig(*configfile)

	mux := http.NewServeMux()
	for _, c := range config {
		pattern, handler, err := c.ToHandler()
		if err != nil {
			fmt.Printf("Error creating handler: %v\n", err)
			os.Exit(1)
		}
		log.Info().Msgf("Serving %s on %s", c.Directory, pattern)
		mux.HandleFunc(pattern, handler)
	}

	httpserver := http.Server{
		Addr:    *addr,
		Handler: mux,
	}
	log.Info().Msg("Listening on " + *addr)
	httpserver.ListenAndServe()
}

func loadConfig(configfile string) []server.Config {
	var config []server.Config
	if configfile != "" {
		file, err := os.Open(configfile)
		if err != nil {
			fmt.Printf("Error opening config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			fmt.Printf("Error decoding config file: %v\n", err)
			os.Exit(1)
		}
	}
	return config
}
