package utils

import (
	"math"
	"math/rand/v2"
	"time"
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
