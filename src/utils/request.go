package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"savannahtech/src/log"
	"savannahtech/src/types"
)

/*
Exponential backoff algorithm

This algorithm calculates the wait time between retries based on the number of retries and a maximum backoff time.

The formula is:

wait_time = min(2^n + random_number_milliseconds, maximum_backoff)

where:
n is the number of retries
random_number_milliseconds is a random number between 0 and 50000
maximum_backoff is the maximum backoff time in milliseconds
*/
func ExponentialBackoff(n uint, maximun_backoff float64) time.Duration {
	// Generate a random number of milliseconds up to 1000
	random_number_milliseconds := rand.Float64() * 200000

	// Calculate the wait time
	var wait_time float64 = math.Min((math.Exp2(float64(n)) + random_number_milliseconds), maximun_backoff)

	return time.Duration(wait_time) * time.Millisecond
}

/*
Helper function to get the next page URL from the "Link" header

Example of "Link" header: <https://api.github.com/repositories/1/commits?page=2>; rel="next", <https://api.github.com/repositories/1/commits?page=3>; rel="last"
*/
func GetNextPageURL(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

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

// Helper function to make an HTTP request with exponential backoff
func fetchWithBackoff(url string, maxRetries int, maximumBackoff float64) (*http.Response, error) {
	log.InfoLogger.Println("Fetching data from", url)

	var resp *http.Response
	var err error

	for i := 0; i < maxRetries; i++ {

		resp, err = http.Get(url)

		if err != nil || resp.StatusCode != http.StatusOK {
			log.ErrorLogger.Println("Error:", err)

			if err != nil {
				log.InfoLogger.Println("Error:", err)
			} else {
				switch resp.StatusCode {
				case http.StatusNotFound:
					return nil, fmt.Errorf("repository not found")
				default:
					log.InfoLogger.Println("Attempt", i+1, "fetching data from:", url, "failed with status code:", resp.StatusCode)
				}
			}

			if i < maxRetries-1 {
				duration := ExponentialBackoff(uint(i), maximumBackoff)
				log.InfoLogger.Println("Sleeping for", duration)
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

// FetchCommitsChan fetches commits from a given URL
func FetchCommits(url string, commitsChan chan<- []types.Commit) error {
	maximumBackoff := 3200 * 1000.0
	maxRetries := 10

	for {
		resp, err := fetchWithBackoff(url, maxRetries, maximumBackoff)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data := []types.Commit{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		commitsChan <- data

		nextURL := GetNextPageURL(resp.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		url = nextURL
	}

	return nil
}

// FetchRepository fetches a repository from a given URL
func FetchRepository(url string) (types.Repository, error) {
	maximumBackoff := 3200 * 1000.0
	maxRetries := 10

	resp, err := fetchWithBackoff(url, maxRetries, maximumBackoff)
	if err != nil {
		if err.Error() == "404 Not Found" {
			return types.Repository{}, fmt.Errorf("repository not found")
		}
		return types.Repository{}, err
	}
	defer resp.Body.Close()

	var repository types.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return types.Repository{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return repository, nil
}
