package watsonxai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
	"wxai-rag/configs"
)

type Client struct {
	APIKey      string
	ProjectID   string
	APIEndpoint string
	AccessToken string
	HTTPClient  *http.Client
	TokenExpiry time.Time
	tokenMutex  sync.Mutex
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TextGenerationPayload struct {
	ModelID    string                 `json:"model_id"`
	Input      string                 `json:"input"`
	Parameters map[string]interface{} `json:"parameters"`
	ProjectID  string                 `json:"project_id"`
}

type TextGenerationResult struct {
	GeneratedText string `json:"generated_text"`
}

type TextGenerationResponse struct {
	Results []TextGenerationResult `json:"results"`
}

func NewClient(cfg configs.WatsonXAIConfig) (*Client, error) {
	client := &Client{
		APIKey:      cfg.APIKey,
		ProjectID:   cfg.ProjectID,
		APIEndpoint: cfg.APIEndpoint,
		HTTPClient:  &http.Client{},
	}
	// Validate the client / set up initial connections
	if client.APIKey == "" || client.APIEndpoint == "" || client.ProjectID == "" {
		return nil, fmt.Errorf("invalid WatsonX.AI configuration: API key, endpoint, and project id must be set")
	}

	// Get wxai access token
	if err := client.GetAccessToken(); err != nil {
		log.Fatalf("Failed to get access token: %v", err)
	}

	return client, nil
}

func (c *Client) GetAccessToken() error {
	c.tokenMutex.Lock()
	defer c.tokenMutex.Unlock()

	if time.Now().Before(c.TokenExpiry) {
		return nil
	}

	url := "https://iam.cloud.ibm.com/identity/token"
	data := "grant_type=urn:ibm:params:oauth:grant-type:apikey&apikey=" + c.APIKey
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get access token, status code: %d", resp.StatusCode)
	}

	var tokenResp AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("error decoding response: %w", err)
	}

	c.AccessToken = tokenResp.AccessToken
	c.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Add(-1 * time.Minute) // Subtract 1 minute as a buffer
	return nil
}

func (c *Client) ensureValidToken() error {
	if time.Now().After(c.TokenExpiry) {
		return c.GetAccessToken()
	}
	return nil
}

func (c *Client) TextGeneration(payload TextGenerationPayload) (string, error) {
	if err := c.ensureValidToken(); err != nil {
		return "", err
	}

	url := c.APIEndpoint + "/ml/v1/text/generation?version=2023-05-29"
	payload.ProjectID = c.ProjectID

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed, status code: %d, response: %s", resp.StatusCode, body)
	}

	var textResp TextGenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&textResp); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(textResp.Results) > 0 {
		return textResp.Results[0].GeneratedText, nil
	}

	return "", fmt.Errorf("no generated text found in response")
}

func (c *Client) TextGenerationStream(payload TextGenerationPayload) error {
	url := c.APIEndpoint + "/ml/v1/text/generation_stream?version=2023-05-02"
	payload.ProjectID = c.ProjectID

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed, status code: %d, response: %s", resp.StatusCode, body)
	}

	reader := io.Reader(resp.Body)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading stream: %w", err)
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buf[:n]))
	}

	return nil
}
