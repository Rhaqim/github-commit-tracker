package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"savannahtech/internal/log"
	"savannahtech/internal/types"
	"savannahtech/internal/utils"
)

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		n              uint
		maximumBackoff float64
		expectedMax    time.Duration
	}{
		{0, 1000.0, 1000 * time.Millisecond},
		{1, 1000.0, 1000 * time.Millisecond},
		{2, 1000.0, 1000 * time.Millisecond},
		{3, 2000.0, 2000 * time.Millisecond},
		{10, 10000.0, 10000 * time.Millisecond},
	}

	for _, tt := range tests {
		duration := utils.ExponentialBackoff(tt.n, tt.maximumBackoff)
		if duration > tt.expectedMax {
			t.Errorf("expected duration to be less than or equal to %v, got %v", tt.expectedMax, duration)
		}
	}
}

func TestGetNextPageURL(t *testing.T) {
	tests := []struct {
		linkHeader  string
		expectedURL string
	}{
		{"<https://api.github.com/repositories/1/commits?page=2>; rel=\"next\", <https://api.github.com/repositories/1/commits?page=3>; rel=\"last\"", "https://api.github.com/repositories/1/commits?page=2"},
		{"<https://api.github.com/repositories/1/commits?page=3>; rel=\"last\"", ""},
		{"", ""},
	}

	for _, tt := range tests {
		nextURL := utils.GetNextPageURL(tt.linkHeader)
		if nextURL != tt.expectedURL {
			t.Errorf("expected URL %v, got %v", tt.expectedURL, nextURL)
		}
	}
}

func TestFetchCommits(t *testing.T) {
	log.Init(false)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Simulating only one page of results with no pagination
		w.Header().Set("Link", ``)
		fmt.Fprintln(w, `[{"sha":"abc123","author": {"login":"Author"},"message":"Commit message","date":"2023-07-30T12:00:00Z","url":"http://example.com"}]`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := server.URL

	commitChan := make(chan []types.Commit)
	errChan := make(chan error)
	defer close(errChan)

	// Set a timeout for the test
	testTimeout := time.After(2 * time.Second)

	// Run the fetch operation in a goroutine
	go func() {
		err := utils.FetchCommits(url, commitChan)
		if err != nil {
			errChan <- err
		}
		close(commitChan) // Close channel after fetching all commits
	}()

	// Variable to store any potential error from the test
	var testErr error

	// Process the commits from the channel
Loop:
	for {
		select {
		case commit, ok := <-commitChan:
			if !ok {
				break Loop
			}
			if len(commit) != 1 {
				testErr = fmt.Errorf("expected 1 commit, got %d", len(commit))
				break Loop
			}

			if commit[0].Sha != "abc123" {
				testErr = fmt.Errorf("expected SHA to be 'abc123', got %v", commit[0].Sha)
				break Loop
			}
		case err := <-errChan:
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			break Loop
		case <-testTimeout:
			t.Fatalf("test timed out")
			break Loop
		}
	}

	// Check if there was an error during test validation
	if testErr != nil {
		t.Fatalf("test validation failed: %v", testErr)
	}
}

func TestFetchRepository(t *testing.T) {
	log.Init(false)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"name":"example-repo","owner":{"login":"example-owner"}}`)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	url := server.URL
	repo, err := utils.FetchRepository(url)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.Name != "example-repo" {
		t.Errorf("expected repository name 'example-repo', got %v", repo.Name)
	}

	if repo.Owner.Login != "example-owner" {
		t.Errorf("expected repository owner 'example-owner', got %v", repo.Owner.Login)
	}
}
