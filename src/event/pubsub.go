package event

import (
	"context"
	"encoding/json"
	"savannahtech/src/database"
	"savannahtech/src/log"
	"savannahtech/src/types"

	"github.com/google/uuid"
)

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

func (q *EventQueue) Publish(event types.Event) error {
	event.ID = uuid.New().String()

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return database.Redis.Publish(q.ctx, q.Key, eventData).Err()
}

func (q *EventQueue) Subscribe(handler func(event types.Event)) error {
	sub := database.Redis.Subscribe(q.ctx, q.Key)
	ch := sub.Channel()

	for msg := range ch {
		var event types.Event
		err := json.Unmarshal([]byte(msg.Payload), &event)
		if err != nil {
			log.ErrorLogger.Printf("Failed to unmarshal event: %v", err)
			continue
		}

		handler(event)
	}

	return sub.Close()
}
