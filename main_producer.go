package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load configuration from the config file
	err := godotenv.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Start the Kafka producer
	fmt.Println("Starting Kafka producer")
	startProducer()
}
