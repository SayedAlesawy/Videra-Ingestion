package ingest

import (
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// workerOnlineHandler Handler for worker coming online events
func (manager *IngestionManager) workerOnlineHandler(pid int, worker process.Process) {
	log.Println(logPrefix, "New worker came online with pid:", pid)

	//Set worker status to ready
	err := manager.setAsReady(pid)
	if errors.IsError(err) {
		log.Println(logPrefix, "Couldn't set ready status for worker pid:", pid)

		return
	}

	manager.workersListMutex.Lock()

	//Add newly joined worker
	manager.workers[pid] = worker

	manager.workersListMutex.Unlock()
}
