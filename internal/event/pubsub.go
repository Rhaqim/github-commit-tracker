package event

import (
	"context"
	"encoding/json"
	"savannahtech/src/database"
	"savannahtech/src/log"
	"savannahtech/src/types"

	"github.com/google/uuid"
)

/*
EventQueue is a queue for events.

It has a key associated with it, which is used to publish and subscribe to events.
*/
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

/*
Publish publishes an event to the event queue.

It marshals the event into JSON format and publishes it to the Redis channel associated with the event queue.
*/
func (q *EventQueue) Publish(event types.Event) error {
	event.ID = uuid.New().String()

	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return database.Redis.Publish(q.ctx, q.Key, eventData).Err()
}

/*
Subscribe subscribes to the event queue and handles events.

It creates a new subscription to the Redis channel associated with the event queue.

It unmarshals the event data from the Redis message payload and passes it to the provided handler function.

If an error occurs during the unmarshaling of the event data, it logs the error and continues to the next message.
*/
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
