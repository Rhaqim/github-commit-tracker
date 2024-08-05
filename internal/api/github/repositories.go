package github

import (
	"fmt"
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/core/entities"
)

// GetRepositoryCommits fetches the commits of a repository.
func (c *Client) GetRepository(owner, repo string) ([]entities.Repository, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits", c.baseURL, owner, repo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var commits []entities.Repository
	if err := decodeResponse(res, &commits); err != nil {
		return nil, err
	}

	return commits, nil
}
