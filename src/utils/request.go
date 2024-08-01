package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand/v2"
	"net/http"
	"strings"
	"time"

	"savannahtech/src/types"
)

// MakeRequest tries to send a request with retries on failure.
// func MakeRequest(url string) ([]byte, error) {
// 	// exponential backoff
// 	var resp *http.Response
// 	var err error

// 	// Define maximum backoff time in milliseconds
// 	maximum_backoff := 32000.0 // 32 seconds

// 	// Define maximum number of retries
// 	maxRetries := 10

// 	for i := 0; i < maxRetries; i++ {
// 		resp, err = http.Get(url)
// 		if err != nil || resp.StatusCode != http.StatusOK {
// 			if err != nil {
// 				log.Println("error", err)
// 			}
// 			log.Println("status code", resp.StatusCode)
// 			log.Println("retrying in", ExponentialBackoff(uint(i), maximum_backoff))
// 			time.Sleep(ExponentialBackoff(uint(i), maximum_backoff))
// 			continue
// 		}
// 		break
// 	}

// 	log.Println("Response status code", resp.StatusCode)
// 	log.Print("\n")

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// check response status
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("failed to get data for url: %s", url)
// 	}

// 	return body, nil
// }

// func FetchData[T any](url string, makeRequest types.RequestFunc[T]) ([]T, error) {

// 	return makeRequest(url)

// }

// // MakeRequest tries to send a request with retries on failure.
// func MakeRequest2[T any](url string) ([]T, error) {
// 	// Exponential backoff settings
// 	var resp *http.Response
// 	var err error
// 	maximum_backoff := 32000.0 // 32 seconds
// 	maxRetries := 10
// 	var allData []T

// 	for {
// 		for i := 0; i < maxRetries; i++ {
// 			resp, err = http.Get(url)
// 			if err != nil || resp.StatusCode != http.StatusOK {
// 				if err != nil {
// 					log.Println("Error:", err)
// 				} else {
// 					log.Println("Status code:", resp.StatusCode)
// 				}
// 				time.Sleep(ExponentialBackoff(uint(i), maximum_backoff))
// 				continue
// 			}
// 			break
// 		}

// 		if err != nil {
// 			return nil, fmt.Errorf("failed to make request: %w", err)
// 		}

// 		body, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read response body: %w", err)
// 		}

// 		var data []T

// 		err = json.Unmarshal(body, &data)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to unmarshal data: %w", err)
// 		}

// 		// Accumulate the data from this page
// 		allData = append(allData, data...)

// 		// Check if there's a next page
// 		nextURL := getNextPageURL(resp.Header.Get("Link"))
// 		if nextURL == "" {
// 			break // No next page, exit the loop
// 		}

// 		url = nextURL // Update the URL for the next request
// 		resp.Body.Close()
// 	}

// 	return allData, nil
// }

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
					log.Println("Error:", err)
				} else {
					log.Println("Status code:", resp.StatusCode)
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
		return types.Repository{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return types.Repository{}, fmt.Errorf("failed to get data for url: %s", url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.Repository{}, err
	}

	var repository types.Repository
	err = json.Unmarshal(body, &repository)
	if err != nil {
		return types.Repository{}, err
	}

	return repository, nil
}
