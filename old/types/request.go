package types

// Define a type for the function signature
type RequestFunc[T any] func(url string) ([]T, error)
