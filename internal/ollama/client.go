// Package ollama provides client functionality for connecting to Ollama
// running on the host machine.
package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is the default Ollama API endpoint.
const DefaultBaseURL = "http://localhost:11434"

// Client handles communication with the Ollama API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Ollama client.
// If baseURL is empty, the default URL (http://localhost:11434) is used.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// IsAvailable checks if Ollama is reachable at the /api/tags endpoint.
// Returns true if Ollama is available, false otherwise.
func (c *Client) IsAvailable() bool {
	url := c.BaseURL + "/api/tags"
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// WaitForAvailable polls Ollama until it becomes available or maxAttempts is reached.
// Returns an error if Ollama is not available after maxAttempts.
func (c *Client) WaitForAvailable(maxAttempts int) error {
	fmt.Printf("Checking Ollama at %s...\n", c.BaseURL)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if c.IsAvailable() {
			fmt.Println("Ollama is available.")
			return nil
		}
		fmt.Printf("Waiting for Ollama... (%d/%d)\n", attempt, maxAttempts)
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("unable to reach Ollama at %s\nEnsure Ollama is running: ollama serve", c.BaseURL)
}

// ModelExists checks if a model exists in Ollama via the /api/show endpoint.
// Returns true if the model exists, false otherwise.
func (c *Client) ModelExists(modelName string) bool {
	url := c.BaseURL + "/api/show"

	reqBody := map[string]string{"model": modelName}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return false
	}

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// EnsureModel checks if a model exists in Ollama and returns an error if not.
// The error message includes a hint to pull the model.
func (c *Client) EnsureModel(modelName string) error {
	fmt.Printf("Checking for model: %s\n", modelName)

	if c.ModelExists(modelName) {
		fmt.Printf("Model %s is available.\n", modelName)
		return nil
	}

	return fmt.Errorf("model %s not found in Ollama\nPull it first: ollama pull %s", modelName, modelName)
}

// GetModels returns a list of models available in Ollama via the /api/tags endpoint.
// This is useful for listing available models.
func (c *Client) GetModels() ([]string, error) {
	url := c.BaseURL + "/api/tags"

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tagsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.Unmarshal(body, &tagsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	models := make([]string, 0, len(tagsResp.Models))
	for _, m := range tagsResp.Models {
		models = append(models, m.Name)
	}

	return models, nil
}
