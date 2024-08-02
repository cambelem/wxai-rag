package watsonxai

import (
	"fmt"
	"wxai-rag/configs"
)

type Client struct {
	APIKey      string
	ProjectID   string
	APIEndpoint string
}

func NewClient(cfg configs.WatsonXAIConfig) (*Client, error) {
	client := &Client{
		APIKey:      cfg.APIKey,
		ProjectID:   cfg.ProjectID,
		APIEndpoint: cfg.APIEndpoint,
	}
	// Validate the client / set up initial connections
	if client.APIKey == "" || client.APIEndpoint == "" || client.ProjectID == "" {
		return nil, fmt.Errorf("invalid WatsonX.AI configuration: API key, endpoint, and project id must be set")
	}

	return client, nil
}

func (c *Client) GenerateText(prompt string) (string, error) {
	return "", nil
}
