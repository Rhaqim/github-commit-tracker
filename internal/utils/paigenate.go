package utils

import "strconv"

func Paigenation(pageStr, sizeStr string) (int, int) {
	// Default values
	page := 1
	size := 10

	// Parse page and size if provided
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
			size = s
		}
	}

	// Calculate offset
	offset := (page - 1) * size

	return size, offset
}
