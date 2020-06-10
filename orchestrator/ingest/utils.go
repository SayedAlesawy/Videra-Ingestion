package ingest

import (
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
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

// correctJob A function to check if the job found in cache matches the job in the event or not
func (manager *IngestionManager) correctJob(jobToken string, jid int64) bool {
	activeJob, exists := manager.getActiveJob(jobToken)
	if !exists {
		log.Println(logPrefix, fmt.Sprintf("Job token: %s doesn't match any job", jobToken))
	}

	job, err := decode(activeJob)
	if errors.IsError(err) {
		log.Println(fmt.Sprintf("%s Unable to decode job: %s, err: %v", logPrefix, activeJob, err))

		return true
	}

	if job.Jid != jid {
		log.Println(logPrefix, fmt.Sprintf("%s Active job cache mismatch, given: %d, found: %d", logPrefix, jid, job.Jid))
		return false
	}

	return true
}
