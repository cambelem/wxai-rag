package configs

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type ElasticsearchConfig struct {
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
}

type WatsonXAIConfig struct {
	APIKey      string `yaml:"api_key"`
	ProjectID   string `yaml:"project_id"`
	APIEndpoint string `yaml:"api_endpoint"`
}

type Config struct {
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	WatsonxAI     WatsonXAIConfig     `yaml:"watsonxai"`
}

func LoadConfig(filepath string) Config {
	var config Config
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}
	return config
}
