package main

import (
	"log"
	"net/http"
	"wxai-rag/api"
	"wxai-rag/configs"
	"wxai-rag/internal/watsonxai"
	"wxai-rag/pkg/elasticsearch"
)

func main() {
	// Load configuration
	config := configs.LoadConfig("configs/config.yaml")

	// Initialize Elasticsearch client
	esClient, err := elasticsearch.NewClient(config.Elasticsearch)
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch client: %v", err)
	}

	// Initialize WatsonX.AI client
	wxClient, err := watsonxai.NewClient(config.WatsonxAI)
	if err != nil {
		log.Fatalf("Failed to initialize WatsonX.AI client: %v", err)
	}

	// Set up API routes
	router := api.SetupRouter(esClient, wxClient)

	// Start the HTTP server
	log.Println("Starting server on :4060")
	if err := http.ListenAndServe(":4060", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
