package ingest

import (
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/health"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/pubsub"
)

// Notify A function to notify the ingestion manager with changes in the process list
func (manager *IngestionManager) Notify(event pubsub.Event) {
	switch event.Type {
	case health.EventProcessOnline:
		manager.workerOnlineHandler(event.Args[0].(int), event.Args[1].(process.Process))

	case health.EventProcessCrashed:
		manager.workerCrashedHandler(event.Args[0].(int))

	case health.EventJobStarted:
		manager.jobStartedHandler(event.Args[0].(int), event.Args[1].(string))

	case health.EventJobCompleted:
		manager.jobCompletedHandler(event.Args[0].(int), event.Args[1].(string))

	default:
		log.Println(logPrefix, "Unknown event type:", event.Type)
	}
}
