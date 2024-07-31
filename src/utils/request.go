package utils

import (
	"context"
	"io"
	"net/http"
	"time"

	"savannahtech/src/types"
)

// MakeRequest tries to send a request with retries on failure.
func MakeRequest(url string) (*http.Response, error) {
	// Set the number of retries and the delay between them
	const maxRetries = 3
	const retryDelay = 2 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		// Set the timeout duration
		timeout := 5 * time.Second

		// Create a context with the timeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Create a new request
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(retryDelay) // Wait before retrying
			continue
		}

		// If successful, return the response
		return resp, nil
	}

	// Return the last error encountered
	return nil, lastErr
}

// func MakeRequest(url string) (*http.Response, error) {
// 	// Set the timeout duration
// 	timeout := 5 * time.Second

// 	// Create a context with the timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()

// 	// Create a new request
// 	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Send the request
// 	client := &http.Client{}

// 	return client.Do(req)
// }

func FetchData(url string, makeRequest types.RequestFunc) ([]byte, error) {
	var body []byte

	resp, err := makeRequest(url)
	if err != nil {
		return body, err
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}
