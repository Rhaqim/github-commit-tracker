package github

// client.go
// This file contains the implementation of the GitHub API client.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Rhaqim/savannahtech/internal/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

const (
	maximumBackoff = 3200 * 1000.0
	maxRetries     = 10
)

// Client is the GitHub API client.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new GitHub API client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.github.com/repos/",
	}
}

func decodeResponse(res *http.Response, v interface{}) error {
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(v)
}

func fetchWithBackoff(url string) (*http.Response, error) {
	logger.InfoLogger.Printf("Fetching data from %s\n", url)

	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {

		resp, err = http.Get(url)

		if err != nil || resp.StatusCode != http.StatusOK {

			if err != nil {
				logger.InfoLogger.Println("Error:", err)
			} else {
				switch resp.StatusCode {
				case http.StatusNotFound:
					logger.ErrorLogger.Printf("Repository not found: %s\n", url)
					return nil, fmt.Errorf("repository not found")
				default:
					logger.InfoLogger.Printf("Attempt %d fetching data from %s failed with status code: %d\n", i+1, url, resp.StatusCode)
				}
			}

			if i < maxRetries-1 {
				duration := utils.ExponentialBackoff(uint(i), maximumBackoff)
				logger.InfoLogger.Printf("Retrying in %s\n", duration)
				time.Sleep(duration)
			}
			continue
		}

		return resp, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
