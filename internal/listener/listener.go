package listener

import (
	"github.com/Rhaqim/savannahtech/internal/events"
)

func StartCommitEventListener() {
	events.StartEventListener(ProcessFunc)
}

func StartEventListeners() {
	StartCommitEventListener()
}
