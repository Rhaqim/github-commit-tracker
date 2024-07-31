package test

import (
	"encoding/json"
	"errors"
	"savannahtech/src/core"
	"savannahtech/src/types"
	"testing"
)

type Commit struct {
	Commit CommitDetail `json:"commit"`
}

type CommitDetail struct {
	Committer Committer `json:"committer"`
	Message   string    `json:"message"`
	Url       string    `json:"url"`
}

type Committer struct {
	Name  string `json:"name"`
	Date  string `json:"date"`
	Email string `json:"email"`
}

func TestFetchAndProcessDataSuccess(t *testing.T) {
	// Mock fetch function
	mockFetchFunc := func(owner, repo string, makeRequest types.RequestFunc) ([]byte, error) {
		mockData := []Commit{
			{
				Commit: CommitDetail{
					Committer: Committer{Name: "Test Committer", Date: "2024-07-30", Email: "test@example.com"},
					Message:   "Test commit message",
					Url:       "https://example.com/commit",
				},
			},
		}
		return json.Marshal(mockData)
	}

	// Mock process function
	mockProcessFunc := func(commit Commit) error {
		if commit.Commit.Message != "Test commit message" {
			t.Errorf("unexpected commit message: %s", commit.Commit.Message)
		}
		return nil
	}

	err := core.FetchAndProcessData("testowner", "testrepo", mockFetchFunc, mockProcessFunc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestFetchAndProcessDataFetchError(t *testing.T) {
	// Mock fetch function that returns an error
	mockFetchFunc := func(owner, repo string, makeRequest types.RequestFunc) ([]byte, error) {
		return nil, errors.New("fetch error")
	}

	// Mock process function
	mockProcessFunc := func(commit Commit) error {
		return nil
	}

	err := core.FetchAndProcessData("testowner", "testrepo", mockFetchFunc, mockProcessFunc)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if err.Error() != "fetch error" {
		t.Errorf("expected fetch error, got %v", err)
	}
}

func TestFetchAndProcessDataProcessError(t *testing.T) {
	// Mock fetch function that returns valid data
	mockFetchFunc := func(owner, repo string, makeRequest types.RequestFunc) ([]byte, error) {
		mockData := []Commit{
			{
				Commit: CommitDetail{
					Committer: Committer{Name: "Test Committer", Date: "2024-07-30", Email: "test@example.com"},
					Message:   "Test commit message",
					Url:       "https://example.com/commit",
				},
			},
		}
		return json.Marshal(mockData)
	}

	// Mock process function that returns an error
	mockProcessFunc := func(commit Commit) error {
		return errors.New("process error")
	}

	err := core.FetchAndProcessData("testowner", "testrepo", mockFetchFunc, mockProcessFunc)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if err.Error() != "process error" {
		t.Errorf("expected process error, got %v", err)
	}
}
