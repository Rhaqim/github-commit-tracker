package github

// client.go
// This file contains the implementation of the GitHub API client.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the GitHub API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new GitHub API client.
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func decodeResponse(res *http.Response, v interface{}) error {
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(v)
}
