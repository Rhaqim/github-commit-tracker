package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/api/github/types"
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

func FetchRepository(url string) (types.Repository, error) {
	resp, err := fetchWithBackoff(url)
	if err != nil {
		if err.Error() == "404 Not Found" {
			return types.Repository{}, fmt.Errorf("repository not found")
		}
		return types.Repository{}, err
	}
	defer resp.Body.Close()

	var repository types.Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return types.Repository{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return repository, nil
}
