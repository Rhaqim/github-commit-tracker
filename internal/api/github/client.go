package github

// client.go
// This file contains the implementation of the GitHub API client.

import (
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
					logger.ErrorLogger.Printf("Attempt %d fetching data from %s failed with status code: %d\n\n", i+1, url, resp.StatusCode)
				}
			}

			if i < maxRetries-1 {
				duration := utils.ExponentialBackoff(uint(i), maximumBackoff)
				logger.ErrorLogger.Printf("Retrying in %s\n", duration)
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
