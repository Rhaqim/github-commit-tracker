package test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"savannahtech/src/core"
	"savannahtech/src/utils"
	"testing"
)

// Mock version of MakeRequest
func mockMakeRequest[T any](url string) ([]T, error) {
	// Check for a specific test URL to simulate different responses
	if url == "https://api.github.com/repos/testowner/testrepo/commits" {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusOK)
		recorder.WriteString(`[{"sha": "abc123"}]`)

		// Return the response body and any error
		result := recorder.Result()

		body, err := io.ReadAll(result.Body)
		if err != nil {
			return nil, err
		}

		var data []T

		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	if url == "https://api.github.com/repos/testowner/testrepo" {
		recorder := httptest.NewRecorder()
		recorder.WriteHeader(http.StatusOK)
		recorder.WriteString(`{"id": "12345"}`)

		result := recorder.Result()

		body, err := io.ReadAll(result.Body)
		if err != nil {
			return nil, err
		}

		var data []T

		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return nil, errors.New("error making request")
}

func TestFetchCommit(t *testing.T) {
	body, err := core.FetchCommit("testowner", "testrepo", mockMakeRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `abc123`
	if body[0].Sha != expected {
		t.Fatalf("expected %s, got %s", expected, body[0].Sha)
	}
}

func TestFetchRepo(t *testing.T) {
	body, err := core.FetchRepo("testowner", "testrepo", mockMakeRequest)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `{"id": "12345"}`
	if string(body) != expected {
		t.Fatalf("expected %s, got %s", expected, string(body))
	}
}

func TestFetchDataError(t *testing.T) {
	failingMakeRequest := func(url string) ([]byte, error) {
		return nil, errors.New("failed to make request")
	}

	body, err := utils.FetchData("https://invalid.url", failingMakeRequest)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if len(body) != 0 {
		t.Fatalf("expected empty body, got %v", string(body))
	}
}
