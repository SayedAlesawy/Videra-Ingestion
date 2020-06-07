package ingest

import (
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/health"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
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
	manager.jobsListMutex.Lock()
	manager.workersListMutex.Lock()
	manager.activeJobsMutex.Lock()

	//Validate that worker exist
	_, workerExists := manager.workers[pid]
	if !workerExists {
		log.Println(logPrefix, fmt.Sprintf("Unknown worker pid: %d", pid))
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d crashed", pid))

	//Check if the crashed worker was executing any jobs
	job, exists := manager.activeJobs[pid]
	if exists {
		//Move the crashed worker job to todo to be retry-ied
		manager.unmarkAsInFlight(job.Jid)
		manager.jobsList[job.Jid] = job

		//Remove it from active jobs
		delete(manager.activeJobs, pid)
	}

	//Remove the worker from the list of online workers and clean up tcp connection
	delete(manager.workers, pid)

	manager.jobsListMutex.Unlock()
	manager.workersListMutex.Unlock()
	manager.activeJobsMutex.Unlock()
}

// jobStartedHandler Handler for jobs starting events
func (manager *IngestionManager) jobStartedHandler(pid int, jid int64) {
	manager.jobsListMutex.Lock()
	manager.workersListMutex.Lock()
	manager.activeJobsMutex.Lock()

	//Validate that worker and job exist
	_, workerExists := manager.workers[pid]
	_, jobExists := manager.jobsList[jid]
	if !workerExists || !jobExists {
		log.Println(logPrefix, fmt.Sprintf("Unknown worker pid: %d or job jid: %d", pid, jid))
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d started executing job with jid: %d", pid, jid))

	//Insert it in active jobs
	manager.activeJobs[pid] = manager.jobsList[jid]

	//Remove it from todo jobs
	manager.unmarkAsInFlight(jid)
	delete(manager.jobsList, jid)

	//Mark worker as busy
	worker := manager.workers[pid]
	worker.Utilization.Busy = true
	manager.workers[pid] = worker

	manager.jobsListMutex.Unlock()
	manager.workersListMutex.Unlock()
	manager.activeJobsMutex.Unlock()
}

// jobCompletedHandler Handler for jobs completion events
func (manager *IngestionManager) jobCompletedHandler(pid int, jid int64) {
	manager.jobsListMutex.Lock()
	manager.workersListMutex.Lock()
	manager.activeJobsMutex.Lock()

	//Validate that worker and job exist
	_, workerExists := manager.workers[pid]
	_, jobExists := manager.activeJobs[pid]
	if !workerExists || !jobExists {
		log.Println(logPrefix, fmt.Sprintf("Unknown worker pid: %d or job jid: %d", pid, jid))
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d completed executing job with jid: %d", pid, jid))

	//Remove it from active jobs
	delete(manager.activeJobs, pid)

	//Mark the worker as not busy
	worker := manager.workers[pid]
	worker.Utilization.Busy = false
	manager.workers[pid] = worker

	manager.jobsListMutex.Unlock()
	manager.workersListMutex.Unlock()
	manager.activeJobsMutex.Unlock()
}
