package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/types"
	"github.com/Rhaqim/savannahtech/old/utils"
)

// GetRepositoryCommit fetches a commit of a repository.
func (c *Client) GetRepositoryCommit(ctx context.Context, owner, repo, sha string) (*entities.Commit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s", c.baseURL, owner, repo, sha)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var commit entities.Commit
	if err := decodeResponse(res, &commit); err != nil {
		return nil, err
	}

	return &commit, nil
}

func FetchCommits(url string, commitsChan chan<- []types.Commit) error {

	for {
		resp, err := fetchWithBackoff(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		data := []types.Commit{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		commitsChan <- data

		nextURL := utils.GetNextPageURL(resp.Header.Get("Link"))
		if nextURL == "" {
			break
		}

		url = nextURL
	}

	return nil
}
