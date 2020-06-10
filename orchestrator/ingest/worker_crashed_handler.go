package ingest

import (
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// workerCrashedHandler Handler for worker crashing events
func (manager *IngestionManager) workerCrashedHandler(pid int) {
	//Validate that worker and job exist
	if !manager.workerExists(pid) {
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d crashed", pid))

	//Check if the worker has active jobs
	activeJobToken, hasJob := manager.hasActiveJob(pid)
	if hasJob {
		//Check in-progress queue, if found, then move to todo
		inProgress, err := manager.findJobInQueue(manager.Queues.InProgress, activeJobToken)
		if errors.IsError(err) {
			log.Println(logPrefix, fmt.Sprintf("%s Error while searching for job in %s for worker pid: %d ",
				logPrefix, manager.Queues.Done, pid))
		} else {
			if inProgress {
				//Move to todo
				err := manager.moveInQueues(activeJobToken, manager.Queues.InProgress, manager.Queues.Todo)
				errors.HandleError(err, fmt.Sprintf("%s Error while moving active job from %s to %s by worker pid: %d",
					logPrefix, manager.Queues.InProgress, manager.Queues.Todo, pid), false)
			} else {
				//Check in done queue
				inDone, err := manager.findJobInQueue(manager.Queues.Done, activeJobToken)
				if errors.IsError(err) {
					log.Println(logPrefix, fmt.Sprintf("Error while searching for job in %s for worker pid: %d ",
						manager.Queues.Done, pid))
				} else {
					if !inDone {
						//Check in todo queue
						inTodo, err := manager.findJobInQueue(manager.Queues.Todo, activeJobToken)
						if errors.IsError(err) {
							log.Println(logPrefix, fmt.Sprintf("Error while searching for job in %s for worker pid: %d ",
								manager.Queues.Done, pid))
						} else {
							if !inTodo {
								//Insert in todo
								err := manager.insertJobsInQueue(manager.Queues.Todo, activeJobToken)
								errors.HandleError(err, fmt.Sprintf("%s Error while inserting job in %s for worker pid: %d",
									logPrefix, manager.Queues.Todo, pid), false)
							}
						}
					}
				}
			}
		}

		//Remove the worker's active job
		err = manager.removeActiveJob(pid)
		errors.HandleError(err, fmt.Sprintf("%s Error while removing active job for worker pid: %d ", logPrefix, pid), false)
	}

	//Remove the worker from the list of online workers and clean up tcp connection
	manager.workersListMutex.Lock()
	delete(manager.workers, pid)
	manager.workersListMutex.Unlock()
}
