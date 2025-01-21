package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	// Load configuration from the config file
	err := godotenv.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration file: %v", err)
	}

	// Create a channel to handle OS signals (graceful shutdown)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// Start Kafka consumers with parallel workers
	fmt.Println("Starting Kafka consumers")
	go startConsumers(3)

	// Wait for termination signal
	<-stopChan
	log.Println("Consumer application stopped.")
}
