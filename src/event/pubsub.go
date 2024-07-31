package event

import (
	"context"
	"encoding/json"
	"log"
	"savannahtech/src/database"

	"github.com/google/uuid"
)

type EventType string

const (
	CommitEvent EventType = "commit"
	RepoEvent   EventType = "repo"
)

type Event struct {
	ID      string    `json:"id"`
	Message string    `json:"message"`
	Type    EventType `json:"type"`
	Owner   string    `json:"owner"`
	Repo    string    `json:"repo"`
}

type EventQueue struct {
	Key string
	ctx context.Context
}

func NewEventQueue(key string) *EventQueue {
	return &EventQueue{
		Key: key,
		ctx: context.Background(),
	}
}

func (q *EventQueue) Publish(event Event) error {
	event.ID = uuid.New().String()

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return database.Redis.Publish(q.ctx, q.Key, eventData).Err()
}

func (q *EventQueue) Subscribe(handler func(event Event)) error {
	sub := database.Redis.Subscribe(q.ctx, q.Key)
	ch := sub.Channel()

	for msg := range ch {
		var event Event
		err := json.Unmarshal([]byte(msg.Payload), &event)
		if err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		handler(event)
	}

	return sub.Close()
}
