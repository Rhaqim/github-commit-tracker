package utils

import (
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"time"

	"savannahtech/src/types"
)

// MakeRequest tries to send a request with retries on failure.
func MakeRequest(url string) ([]byte, error) {
	// exponential backoff
	var resp *http.Response
	var err error

	// Define maximum backoff time in milliseconds
	maximum_backoff := 32000.0 // 32 seconds

	// Define maximum number of retries
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		resp, err = http.Get(url)
		if err != nil || resp.StatusCode != http.StatusOK {
			if err != nil {
				fmt.Println("error", err)
			}
			fmt.Println("status code", resp.StatusCode)
			fmt.Println("retrying in", ExponentialBackoff(uint(i), maximum_backoff))
			time.Sleep(ExponentialBackoff(uint(i), maximum_backoff))
			continue
		}
		break
	}

	fmt.Println("Response status code", resp.StatusCode)
	fmt.Print("\n")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get potentials for user: %s", string(body))
	}

	return body, nil
}

func ExponentialBackoff(n uint, maximun_backoff float64) time.Duration {
	// Generate a random number of milliseconds up to 1000
	random_number_milliseconds := rand.Float64() * 1000

	// Calculate the wait time
	var wait_time float64 = math.Min((math.Exp2(float64(n)) + random_number_milliseconds), maximun_backoff)

	return time.Duration(wait_time) * time.Millisecond
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

// req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Send the request
// 	client := &http.Client{}

// 	return client.Do(req)
// }

func FetchData(url string, makeRequest types.RequestFunc) ([]byte, error) {

	return makeRequest(url)

}
