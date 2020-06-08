package ingest

import (
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/health"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
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
		manager.jobStartedHandler(event.Args[0].(int), event.Args[1].(int64))

	case health.EventJobCompleted:
		manager.jobCompletedHandler(event.Args[0].(int), event.Args[1].(int64))

	default:
		log.Println(logPrefix, "Unknown event type:", event.Type)
	}
}

// workerOnlineHandler Handler for worker coming online events
func (manager *IngestionManager) workerOnlineHandler(pid int, worker process.Process) {
	log.Println(logPrefix, "New worker came online with pid:", pid)

	manager.workersListMutex.Lock()

	//Add newly joined worker
	manager.workers[pid] = worker

	manager.workersListMutex.Unlock()
}

// workerCrashedHandler Handler for worker crashing events
func (manager *IngestionManager) workerCrashedHandler(pid int) {
	//Validate that worker and job exist
	if !manager.workerExists(pid) {
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d crashed", pid))

	//Check if the worker has active jobs
	activeJob, hasJob := manager.hasActiveJob(pid)
	if hasJob {
		//Move jobs from in-progress to todo
		err := manager.moveInQueues(activeJob, manager.Queues.InProgress, manager.Queues.Todo)
		errors.HandleError(err, fmt.Sprintf("%s Error while moving active job from %s to %s by worker pid: %d",
			logPrefix, manager.Queues.InProgress, manager.Queues.Todo, pid), false)

		//Remove the worker's active job
		err = manager.removeActiveJob(pid)
		errors.HandleError(err, fmt.Sprintf("%s Error while removing active job for worker pid: %d ", logPrefix, pid), false)
	}

	//Remove the worker from the list of online workers and clean up tcp connection
	manager.workersListMutex.Lock()
	delete(manager.workers, pid)
	manager.workersListMutex.Unlock()
}

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

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d completed executing job with jid: %d", pid, jid))

	//Remove the active job from in-progress to done
	err := manager.moveInQueues(activeJob, manager.Queues.InProgress, manager.Queues.Done)
	errors.HandleError(err, fmt.Sprintf("%s Error while moving job jid: %d from %s to %s by worker pid: %d",
		logPrefix, jid, manager.Queues.InProgress, manager.Queues.Done, pid), false)

	//Remove it from active jobs
	err = manager.removeActiveJob(pid)
	errors.HandleError(err, fmt.Sprintf("%s Error while removing active job jid:%d for worker pid: %d ", logPrefix, jid, pid), false)

	//Mark the worker as not busy
	manager.updateWorkerBusyStatus(pid, false)
}

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
	cacheKey, jobExists := manager.getActiveJob(pid)
	if !jobExists {
		log.Println(logPrefix, fmt.Sprintf("%s No active jobs for worker pid: %d found in cache", logPrefix, pid))
		return "", false
	}

	return cacheKey, true
}

// correctJob A function to check if the job found in cache matches the job in the event or not
func (manager *IngestionManager) correctJob(cacheKey string, jid int64) bool {
	job, err := decode(cacheKey)
	if errors.IsError(err) {
		log.Println(fmt.Sprintf("%s Unable to decode job: %s, err: %v", logPrefix, cacheKey, err))

		return true
	}

	if job.Jid != jid {
		log.Println(logPrefix, fmt.Sprintf("%s Active job cache mismatch, given: %d, found: %d", logPrefix, jid, job.Jid))
		return false
	}

	return true
}
