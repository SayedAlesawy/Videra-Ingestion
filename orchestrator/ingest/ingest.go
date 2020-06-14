package ingest

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/redis"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/params"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Ingestion-Manager]"

// ingestionManagerOnce Used to garauntee thread safety for singleton instances
var ingestionManagerOnce sync.Once

// ingestionManagerInstance A singleton instance of the ingestion manager object
var ingestionManagerInstance *IngestionManager

// IngestionManagerInstance A function to return an ingestion manager instance
func IngestionManagerInstance() *IngestionManager {
	ingestionManagerOnce.Do(func() {
		configManager := config.ConfigurationManagerInstance("config/config_files")
		configObj := configManager.IngestionManagerConfig("ingestion_manager.yaml")
		params := params.OrchestratorParamsInstance()

		cacheInstance, err := redis.Instance(configManager.RedisConfig())
		errors.HandleError(err, fmt.Sprintf("%s Unable to connect to caching layer", logPrefix), true)

		manager := IngestionManager{
			startIdx:          params.StartIdx,
			frameCount:        params.FrameCount,
			jobSize:           configObj.JobSize,
			workers:           make(map[int]process.Process),
			workersListMutex:  &sync.Mutex{},
			cache:             cacheInstance,
			cachePrefix:       params.ExecutionGroupID,
			checkDoneInterval: time.Duration(configObj.CheckDoneInterval) * time.Second,
		}

		manager.getQueueNames(configObj.Queues)

		ingestionManagerInstance = &manager
	})

	return ingestionManagerInstance
}

// Start Starts the ingestion manager job scheduling routine
func (manager *IngestionManager) Start(done chan os.Signal) {
	log.Println(logPrefix, "Starting Ingestion Manager")

	log.Println(logPrefix, fmt.Sprintf("Inserting jobs in %s", manager.queues.Todo))

	manager.jobCount = manager.populateJobsPool()

	log.Println(logPrefix, fmt.Sprintf("Successfully inserted %d jobs in %s", manager.jobCount, manager.queues.Todo))

	//Check when ingestion is done
	go manager.checkDone(done)
}

// Shutdown A function to shutdown the ingestion manager
func (manager *IngestionManager) Shutdown() {
	log.Println(logPrefix, "Ingestion Manager cleaning cache")

	manager.flushCache()

	log.Println(logPrefix, "Ingestion Manager shutdown successfully")
}

// checkDone A function to check if all ingestion jobs are done or not
func (manager *IngestionManager) checkDone(done chan os.Signal) {
	for range time.Tick(manager.checkDoneInterval) {
		doneLen, err := manager.getQueueLength(manager.queues.Done)
		if errors.IsError(err) {
			log.Println(logPrefix, fmt.Sprintf("Unable to check %s length, err: %v", manager.queues.Done, err))
			continue
		}

		//If all jobs are done, then terminate
		if int(doneLen) == manager.jobCount {
			log.Println(logPrefix, fmt.Sprintf("Done ingesting all %d jobs. Terminating.", manager.jobCount))

			manager.Shutdown()
			done <- os.Kill

			return
		}
	}
}

// getQueueNames A function to contruct ingestion queue names
func (manager *IngestionManager) getQueueNames(queues config.Queue) {
	manager.queues = Queue{
		Todo:       fmt.Sprintf("%s:%s", manager.cachePrefix, queues.Todo),
		InProgress: fmt.Sprintf("%s:%s", manager.cachePrefix, queues.InProgress),
		Done:       fmt.Sprintf("%s:%s", manager.cachePrefix, queues.Done),
	}
}

// populateJobsPool Populates the jobs pool of the ingestion manager
func (manager *IngestionManager) populateJobsPool() int {
	var jobs []interface{}
	jobTokens := make(map[string]string)

	jobsCount := manager.frameCount / manager.jobSize
	remainder := manager.frameCount % manager.jobSize

	if remainder != 0 {
		jobsCount++
	}

	for i, start := int64(0), manager.startIdx; i < jobsCount; i, start = i+1, start+manager.jobSize {
		jobSize := manager.jobSize
		jid := i + 1

		if remainder != 0 && i == jobsCount-1 {
			jobSize = remainder
		}

		job := newIngestionJob(jid, start, jobSize)
		encodedJob, err := job.encode()
		errors.HandleError(err, fmt.Sprintf("%s Unable to encode job: %+v", logPrefix, job), true)

		jobs = append(jobs, jid)
		jobTokens[fmt.Sprintf("%d", jid)] = encodedJob
	}

	err := manager.insertJobsInQueue(manager.queues.Todo, jobs...)
	errors.HandleError(err, fmt.Sprintf("%s Unable to insert todo jobs on start up", logPrefix), true)

	err = manager.insertJobTokens(jobTokens)
	errors.HandleError(err, fmt.Sprintf("%s Unable to insert job keys on start up", logPrefix), true)

	return len(jobs)
}
