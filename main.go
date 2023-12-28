package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/harr1424/Go-Notify/gonotify"
)

var svc *dynamodb.Client

func main() {

	// Load AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(gonotify.Region))
	if err != nil {
		log.Fatal("Error loading SDK config: ", err)
	}

	// Create DynamoDB client and provide it to gonotify package methods
	svc = dynamodb.NewFromConfig(cfg)
	gonotify.InitializeDynamoDBClient(svc)

	// Attempt to read existing data from DynamoDB
	// Limit to three attempts 5 seconds apart
	maxRetries := 3
	retryDelay := 5 * time.Second
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := gonotify.ReadRemoteTableContents()
		if err == nil {
			break // Success
		}

		log.Printf("Attempt %d: Failed to read remote table contents: %v", attempt, err)

		if attempt < maxRetries {
			log.Printf("Retrying in %v...", retryDelay)
			time.Sleep(retryDelay)
		}

		log.Fatal("Unable to load remote table contents: ", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", gonotify.RegisterToken)
	mux.HandleFunc("/add_location", gonotify.HandleLocationAdd)
	mux.HandleFunc("/remove_location", gonotify.HandleLocationRemove)

	// Every 12 hours check all saved locations for frost
	go func() {
		for {
			gonotify.CheckAllLocationsForFrost()
			time.Sleep(12 * time.Hour)
		}
	}()

	log.Fatal(http.ListenAndServe("0.0.0.0:5050", mux))
	//log.Fatal(http.ListenAndServeTLS(":5050", "localhost.crt", "localhost.key", nil)) // support TLS when available
}
