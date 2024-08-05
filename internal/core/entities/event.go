package entities

type EventType string

const (
	NewRepo     EventType = "new-repo"
	CommitEvent EventType = "commit"
	RepoEvent   EventType = "repo"
)

type Event struct {
	ID      string    `json:"id"`
	From    string    `json:"from"`
	Message string    `json:"message"`
	Type    EventType `json:"type"`
	Owner   string    `json:"owner"`
	Repo    string    `json:"repo"`
}
