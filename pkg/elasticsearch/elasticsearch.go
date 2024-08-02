package elasticsearch

import (
	"fmt"
	"wxai-rag/configs"

	"github.com/elastic/go-elasticsearch/v8"
)

type Client struct {
	es *elasticsearch.Client
}

func NewClient(cfg configs.ElasticsearchConfig) (*Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating the Elasticsearch client: %w", err)
	}

	// Test the connection
	res, err := es.Info()
	if err != nil {
		return nil, fmt.Errorf("error getting response from Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error: %s", res.String())
	}

	return &Client{es: es}, nil
}

func (c *Client) Search(index string, query map[string]interface{}) (interface{}, error) {
	return "", nil
}
