package github

import (
	"encoding/json"
	"fmt"

	"github.com/Rhaqim/savannahtech/internal/core/types"
	"github.com/Rhaqim/savannahtech/internal/utils"
)

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
