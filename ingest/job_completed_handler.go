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
	isDone, exists := manager.doneJobSet[jid]

	if exists && isDone {
		inDoneHardCheck, err := manager.findJobInQueue(manager.queues.Done, jid)

		if errors.IsError(err) {
			log.Println(logPrefix, fmt.Sprintf("Error while searching for job jid: %s in %s for worker pid: %d ",
				jid, manager.queues.Done, pid))
		}
		if inDoneHardCheck {
			manager.doneJobSet[jid] = true // mark job as done

			log.Println(logPrefix, fmt.Sprintf("Job jid: %s is in %s for worker pid: %d. rejecting done signal",
				jid, manager.queues.Done, pid))

			err = manager.removeJobFromeQueue(manager.queues.InProgress, jid)
			errors.HandleError(err, fmt.Sprintf("%s Error while inserting job jid: %s in %s for worker pid: %d",
				logPrefix, jid, manager.queues.Todo, pid), false)

			manager.setAsReady(pid)
		} else {
			manager.markJobAsDone(jid, pid)
		}

	} else {
		manager.markJobAsDone(jid, pid)
	}
}

func (manager *IngestionManager) markJobAsDone(jid string, pid int) {
	manager.moveInQueues(jid, manager.queues.InProgress, manager.queues.Done)
	//Mark the worker as not busy
	manager.updateWorkerBusyStatus(pid, false)
	// Allow worker to take a new job
	manager.setAsReady(pid)

	manager.doneJobSet[jid] = true // mark job as done
	log.Println(logPrefix, fmt.Sprintf("Worker with pid: %d completed executing job with jid: %s", pid, jid))

	log.Println(logPrefix, fmt.Sprintf("checking jid: %s for subsquent jobs", jid))
	manager.addNextJobInPipeline(jid)
}

// addNextJobInPipeline Checks if the job is the final stage in its pipeline
// if not, it fetches the next action to be executed
// and adds it to the todo queue
func (manager *IngestionManager) addNextJobInPipeline(jid string) {
	jidNum, err := strconv.Atoi(jid)
	errors.HandleError(err, fmt.Sprintf("%s Error while inserting subsquent job in pipeline for job: %s",
		logPrefix, jid), false)

	jobData, fetchStatus := manager.getJob(jidNum)
	if fetchStatus == false {
		log.Println(logPrefix, fmt.Sprintf("Failed to fetch next job data"))
		return
	}

	actionType := jobData.Action
	nextAction, exists := actionPipeline[actionType]
	if nextAction != nullAction && exists {
		nextJid := manager.jobCount + 1
		jobData.Action = actionPipeline[actionType]

		encodedJob, err := jobData.encode()
		errors.HandleError(err, fmt.Sprintf("%s Unable to encode job: %+v", logPrefix, jobData), false)

		jobTokens := map[string]string{fmt.Sprintf("%d", nextJid): encodedJob}

		err = manager.insertJobTokens(jobTokens)
		errors.HandleError(err, fmt.Sprintf("%s Error while inserting next job jid: %s in %s",
			logPrefix, jid, manager.queues.Todo), false)

		err = manager.insertJobsInQueue(manager.queues.Todo, nextJid)
		errors.HandleError(err, fmt.Sprintf("%s Error while inserting next job jid: %s in active jobs area",
			logPrefix, jid), false)

		log.Println(logPrefix, "Successfully inserted next job in pipeline with jid: ", nextJid)
		manager.jobCount = nextJid
	}
}
