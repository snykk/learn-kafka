package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func startProducer() {
	// Load Kafka configuration
	config := loadConfig()

	// Initialize Sarama producer
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true

	fmt.Println("config", config)
	producer, err := sarama.NewSyncProducer([]string{config.Kafka.BootstrapServers}, saramaConfig)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// Define the Kafka topic
	topic := config.Kafka.Topic

	// Send 10 messages to Kafka
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("Message-%d", i)

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		})
		if err != nil {
			log.Printf("Failed to send message: %v", err)
		} else {
			log.Printf("Message '%s' sent to partition %d at offset %d", message, partition, offset)
		}

		// Add delay between messages
		time.Sleep(500 * time.Millisecond)
	}
}
