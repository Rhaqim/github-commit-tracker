package entities

type EventType string

const (
	PeriodEvnt  EventType = "Period"
	CommitEvent EventType = "commit"
)

type Event struct {
	Owner     string    `json:"owner"`
	Repo      string    `json:"repo"`
	StartDate string    `json:"startDate"`
	Type      EventType `json:"type"`
}
