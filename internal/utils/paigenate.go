package utils

import (
	"strconv"
	"strings"
)

/*
Paigenation parses the page and size query parameters and returns the size and offset values.

If the page or size is not provided or invalid, it uses the default values (page=1, size=10).
*/
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
