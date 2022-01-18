package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/marcsanmi/query-benchmarking/internal"
	"github.com/marcsanmi/query-benchmarking/internal/parser"
	"github.com/marcsanmi/query-benchmarking/internal/promscale"
	"github.com/marcsanmi/query-benchmarking/internal/reporter"
	"github.com/marcsanmi/query-benchmarking/internal/scheduler"
)

type AppConfig struct {
	Promscale promscale.Config
}

func main() {
	filePath := flag.String("file", "resources/obs-queries.csv", "CSV file path with all the PromQL queries to run")
	workers := flag.Int("workers", 1, "Number of worker processes in mi")
	timeout := flag.Int("timeout", 5000, "Http client timeout in milliseconds.")
	flag.Parse()

	// Process app config.
	var cfg AppConfig
	err := envconfig.Process("APP", &cfg)
	if err != nil {
		log.Fatalf("Unable to process app config: %v", err.Error())
	}

	// Attempt to open the CSV file.
	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Unable to read input csv file: %v", err.Error())
	}
	defer f.Close()

	// Parse csv file into queries.
	queryParser := parser.New(f)
	queries := queryParser.ParseQueries()

	// Give some work to the workers!
	done := make(chan struct{})
	queryChan := make(chan internal.Query)
	go func() {
		for _, q := range queries {
			queryChan <- q
		}

		// Send done signal to the workers.
		for i := 0; i < *workers; i++ {
			done <- struct{}{}
		}
	}()

	// Create basic http client to inject it into the Promscale client.
	httpClient := http.Client{
		Timeout: time.Duration(*timeout * int(time.Millisecond)),
	}

	promscaleClient, err := promscale.NewClient(&httpClient, cfg.Promscale)
	if err != nil {
		log.Fatalf("Failed to create promscale Client: %v", err.Error())
	}

	// Run the scheduler.
	s := scheduler.New(*workers, queryChan, done, promscaleClient, reporter.New())
	s.Run()
}
