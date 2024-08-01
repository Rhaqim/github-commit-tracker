package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"savannahtech/src/log"
	"savannahtech/src/types"
)

func ExponentialBackoff(n uint, maximun_backoff float64) time.Duration {
	// Generate a random number of milliseconds up to 1000
	random_number_milliseconds := rand.Float64() * 1000

	// Calculate the wait time
	var wait_time float64 = math.Min((math.Exp2(float64(n)) + random_number_milliseconds), maximun_backoff)

	return time.Duration(wait_time) * time.Millisecond
}

// Helper function to get the next page URL from the "Link" header
func getNextPageURL(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	// Example of "Link" header: <https://api.github.com/repositories/1/commits?page=2>; rel="next", <https://api.github.com/repositories/1/commits?page=3>; rel="last"
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		parts := strings.Split(strings.TrimSpace(link), ";")

		if len(parts) == 2 && strings.TrimSpace(parts[1]) == `rel="next"` {
			url := strings.Trim(parts[0], "<>")
			return url
		}
	}

	return ""
}

func FetchCommits(url string) ([]types.Commit, error) {
	// Exponential backoff settings
	var resp *http.Response
	var err error

	maximum_backoff := 32000.0 // 32 seconds
	maxRetries := 10

	var allCommits []types.Commit

	for {
		for i := 0; i < maxRetries; i++ {
			resp, err = http.Get(url)

			if err != nil || resp.StatusCode != http.StatusOK {

				if err != nil {
					log.ErrorLogger.Println("Error:", err)
				} else {
					log.InfoLogger.Println("Status code:", resp.StatusCode)
				}

				time.Sleep(ExponentialBackoff(uint(i), maximum_backoff))
				continue
			}

			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to make request: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		var data []types.Commit

		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
		}

		// Accumulate the data from this page
		allCommits = append(allCommits, data...)

		// Check if there's a next page
		nextURL := getNextPageURL(resp.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		// Update the URL for the next request
		url = nextURL
		resp.Body.Close()
	}

	return allCommits, nil
}

func FetchRepository(url string) (types.Repository, error) {
	resp, err := http.Get(url)
	if err != nil {
		return types.Repository{}, fmt.Errorf("error fetching repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Repository{}, fmt.Errorf("failed to get data for url: %s", url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Repository{}, fmt.Errorf("error reading response body: %w", err)
	}

	var repository types.Repository
	err = json.Unmarshal(body, &repository)
	if err != nil {
		return types.Repository{}, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return repository, nil
}
