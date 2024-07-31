package types

import "net/http"

// Define a type for the function signature
type RequestFunc func(url string) (*http.Response, error)
