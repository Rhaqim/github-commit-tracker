package github

import (
	"encoding/json"
	"fmt"

	"github.com/Rhaqim/savannahtech/internal/core/types"
)

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
