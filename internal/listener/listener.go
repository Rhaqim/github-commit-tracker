package listener

import (
	"github.com/Rhaqim/savannahtech/internal/events"
)

func StartEventListeners() {
	events.StartEventListener(ProcessFunc)
}
