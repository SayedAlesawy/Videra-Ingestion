package health

import (
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/pubsub"
)

// Event types sent by the monitor
const (
	EventProcessOnline  = "monitor-event:process-online"
	EventProcessCrashed = "monitor-event:process-crashed"
	EventJobStarted     = "monitor-event:job-started"
	EventJobCompleted   = "monitor-event:job-completed"
)

// notify A function to synchronously notify subscribers with monitor events
func (monitorObj *Monitor) notify(event pubsub.Event) {
	monitorObj.subscribersListMutex.Lock()

	for _, subscriber := range monitorObj.subscribers {
		subscriber.Notify(event)
	}

	monitorObj.subscribersListMutex.Unlock()
}
