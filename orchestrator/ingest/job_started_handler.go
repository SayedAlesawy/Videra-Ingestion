package ingest

import (
	"fmt"
	"log"
)

// jobStartedHandler Handler for jobs starting events
func (manager *IngestionManager) jobStartedHandler(pid int, jid string) {
	//Validate that worker exists
	if !manager.workerExists(pid) {
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d started executing job with jid: %s", pid, jid))

	//Mark worker as busy
	manager.updateWorkerBusyStatus(pid, true)
}
