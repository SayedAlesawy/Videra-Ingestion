package ingest

import (
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// jobCompletedHandler Handler for jobs completion events
func (manager *IngestionManager) jobCompletedHandler(pid int, jid string) {
	//Validate that worker exists
	if !manager.workerExists(pid) {
		return
	}

	//Check that the job is in done queue
	inDone, err := manager.findJobInQueue(manager.queues.Done, jid)
	if errors.IsError(err) {
		log.Println(logPrefix, fmt.Sprintf("Error while searching for job jid: %s in %s for worker pid: %d ",
			jid, manager.queues.Done, pid))
	} else {
		if !inDone {
			//return active job to todo queue for re-try
			log.Println(logPrefix, fmt.Sprintf("Job jid: %s is not in %s for worker pid: %d. Sending back to %s",
				jid, manager.queues.Done, pid, manager.queues.Todo))

			err := manager.insertJobsInQueue(manager.queues.Todo, jid)
			errors.HandleError(err, fmt.Sprintf("%s Error while inserting job jid: %s in %s for worker pid: %d",
				logPrefix, jid, manager.queues.Todo, pid), false)
		}
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d completed executing job with jid: %s", pid, jid))

	//Mark the worker as not busy
	manager.updateWorkerBusyStatus(pid, false)
}
