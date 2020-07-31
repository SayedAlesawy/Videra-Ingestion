package ingest

import (
	"fmt"
	"log"
)

// updateWorkerBusyStatus A function update the worker busy status
func (manager *IngestionManager) updateWorkerBusyStatus(pid int, status bool) {
	manager.workersListMutex.Lock()

	worker := manager.workers[pid]
	worker.Utilization.Busy = status
	manager.workers[pid] = worker

	manager.workersListMutex.Unlock()
}

// workerExists A function to check if a worker exists
func (manager *IngestionManager) workerExists(pid int) bool {
	manager.workersListMutex.Lock()

	_, workerExists := manager.workers[pid]

	manager.workersListMutex.Unlock()

	if !workerExists {
		log.Println(logPrefix, fmt.Sprintf("Unknown worker pid: %d", pid))
	}

	return workerExists
}

// hasActiveJob A function to check if a worker has an active job
func (manager *IngestionManager) hasActiveJob(pid int) (string, bool) {
	jobToken, jobExists := manager.getActiveJobToken(pid)
	if !jobExists {
		log.Println(logPrefix, fmt.Sprintf("%s No active jobs for worker pid: %d found in cache", logPrefix, pid))

		return "", false
	}

	return jobToken, true
}
