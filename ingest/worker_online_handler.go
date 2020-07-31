package ingest

import (
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
)

// workerOnlineHandler Handler for worker coming online events
func (manager *IngestionManager) workerOnlineHandler(pid int, worker process.Process) {
	log.Println(logPrefix, "New worker came online with pid:", pid)

	manager.workersListMutex.Lock()

	//Add newly joined worker
	manager.workers[pid] = worker

	manager.workersListMutex.Unlock()
}
