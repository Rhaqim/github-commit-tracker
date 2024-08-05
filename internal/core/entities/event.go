package entities

type EventType string

const (
	RepoEvent   EventType = "repo"
	CommitEvent EventType = "commit"
	PeriodEvent EventType = "period"
)

// Event represents the event queue data.
type Event struct {
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	StartDate string    `json:"startDate"`
	Type      EventType `json:"type"`
}
