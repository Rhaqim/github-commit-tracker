package events

import "log"

func ProcessEvent(event interface{}) {
	// Process the event based on its type
	log.Println("Processing event:", event)
}
