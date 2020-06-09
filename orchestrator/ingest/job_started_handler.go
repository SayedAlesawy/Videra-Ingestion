package ingest

import (
	"fmt"
	"log"
)

// jobStartedHandler Handler for jobs starting events
func (manager *IngestionManager) jobStartedHandler(pid int, jid int64) {
	//Validate that worker exists
	if !manager.workerExists(pid) {
		return
	}

	//Check if the worker has active jobs
	activeJob, hasJob := manager.hasActiveJob(pid)
	if !hasJob {
		return
	}

	//Check if the job in the healthcheck matches the job in the cache
	correctJob := manager.correctJob(activeJob, jid)
	if !correctJob {
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d started executing job with jid: %d", pid, jid))

	//Mark worker as busy
	manager.updateWorkerBusyStatus(pid, true)
}
