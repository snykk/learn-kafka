package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// KafkaConfig defines the Kafka settings from the configuration file
type KafkaConfig struct {
	BootstrapServers string `yaml:"bootstrap_servers"`
	Topic            string `yaml:"topic"`
	GroupID          string `yaml:"group_id"`
	Retries          int    `yaml:"retries"`
	AutoOffsetReset  string `yaml:"auto_offset_reset"`
}

// Config represents the complete configuration structure
type Config struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

// loadConfig loads Kafka configuration from a YAML file
func loadConfig() *Config {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Failed to open configuration file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to parse configuration file: %v", err)
	}

	return &config
}
