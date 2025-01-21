package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

func startConsumers(workerCount int) {
	// Load Kafka configuration
	config := loadConfig()

	// Initialize Sarama consumer group
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest // Start from the earliest offset

	fmt.Println("config", config)
	consumerGroup, err := sarama.NewConsumerGroup([]string{config.Kafka.BootstrapServers}, config.Kafka.GroupID, saramaConfig)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer consumerGroup.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	// Start consuming messages
	messageChan := make(chan *sarama.ConsumerMessage, 100)
	consumer := Consumer{
		ready:       make(chan bool),
		messageChan: messageChan,
	}

	// Start worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			processMessages(workerID, messageChan, config.Kafka.Retries)
		}(i)
	}

	log.Println("Consumer started listening for messages...")
	for {
		if err := consumerGroup.Consume(ctx, []string{config.Kafka.Topic}, &consumer); err != nil {
			log.Printf("Error during message consumption: %v", err)
		}
		if ctx.Err() != nil {
			break
		}
	}
	close(messageChan)
	wg.Wait()
}

// Consumer implements sarama.ConsumerGroupHandler
type Consumer struct {
	ready       chan bool
	messageChan chan<- *sarama.ConsumerMessage
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.messageChan <- message
		session.MarkMessage(message, "")
	}
	return nil
}

func processMessages(workerID int, messageChan <-chan *sarama.ConsumerMessage, maxRetries int) {
	for msg := range messageChan {
		retries := 0
		for {
			log.Printf("Worker %d processing message: %s", workerID, string(msg.Value))

			// Simulate processing error
			if retries < maxRetries && retries%2 == 0 {
				log.Printf("Worker %d failed to process message, retry %d", workerID, retries)
				retries++
				continue
			}

			log.Printf("Worker %d successfully processed message: %s", workerID, string(msg.Value))
			break
		}
	}
}
