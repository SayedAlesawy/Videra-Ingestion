package ingest

import (
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// jobCompletedHandler Handler for jobs completion events
func (manager *IngestionManager) jobCompletedHandler(pid int, jid int64) {
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

	//Check that the job is in done queue
	inDone, err := manager.findJobInQueue(manager.Queues.Done, activeJob)
	if errors.IsError(err) {
		log.Println(logPrefix, fmt.Sprintf("Error while searching for job jid: %d in %s for worker pid: %d ",
			jid, manager.Queues.Done, pid))
	} else {
		if !inDone {
			//return active job to todo queue for re-try
			log.Println(logPrefix, fmt.Sprintf("Job jid: %d is not in %s for worker pid: %d. Sending back to %s",
				jid, manager.Queues.Done, pid, manager.Queues.Todo))

			err := manager.insertJobsInQueue(manager.Queues.Todo, activeJob)
			errors.HandleError(err, fmt.Sprintf("%s Error while inserting job jid: %d in %s for worker pid: %d",
				logPrefix, jid, manager.Queues.Todo, pid), false)
		}
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d completed executing job with jid: %d", pid, jid))

	//Remove it from active jobs
	err = manager.removeActiveJob(pid)
	errors.HandleError(err, fmt.Sprintf("%s Error while removing active job jid: %d for worker pid: %d ",
		logPrefix, jid, pid), false)

	//Mark the worker as not busy
	manager.updateWorkerBusyStatus(pid, false)
}
