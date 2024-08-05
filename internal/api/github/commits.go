package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
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
