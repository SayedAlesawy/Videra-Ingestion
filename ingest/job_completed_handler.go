package ingest

import (
	"fmt"
	"log"
	"strconv"

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
	log.Println(logPrefix, fmt.Sprintf("checking jid: %s for subsquent jobs", jid))
	manager.addNextJobInPipeline(jid)
}

// addNextJobInPipeline Checks if the job is the final stage in its pipeline
// if not, it fetches the next action to be executed
// and adds it to the todo queue
func (manager *IngestionManager) addNextJobInPipeline(jid string) {
	// ingestion jobs pipeline
	actionPipeline := map[string]string{
		ExecuteAction: MergeAction,
		MergeAction:   NullAction,
		NullAction:    NullAction,
	}

	jidNum, err := strconv.Atoi(jid)
	errors.HandleError(err, fmt.Sprintf("%s Error while inserting subsquent job in pipeline for job: %s",
		logPrefix, jid), false)

	jobData, fetchStatus := manager.getJobToken(jidNum)
	if fetchStatus == false {
		log.Println(logPrefix, fmt.Sprintf("Failed to fetch next job data"))
		return
	}

	actionType := jobData.Action
	if actionPipeline[actionType] != NullAction {
		nextJid := manager.jobCount + 1
		jobData.Action = actionPipeline[actionType]

		encodedJob, err := jobData.encode()
		errors.HandleError(err, fmt.Sprintf("%s Unable to encode job: %+v", logPrefix, jobData), false)

		jobTokens := map[string]string{fmt.Sprintf("%d", nextJid): encodedJob}

		err = manager.insertJobTokens(jobTokens)
		errors.HandleError(err, fmt.Sprintf("%s Error while inserting next job jid: %s in %s",
			logPrefix, jid, manager.queues.Todo), false)

		err = manager.insertJobsInQueue(manager.queues.Todo, nextJid)
		errors.HandleError(err, fmt.Sprintf("%s Error while inserting next job jid: %s in %s",
			logPrefix, jid, manager.queues.Todo), false)

		log.Println(logPrefix, "done adding new job, new target jobs to finish ", nextJid)
		manager.jobCount = nextJid
	}
}
