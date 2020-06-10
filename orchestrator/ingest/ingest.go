package ingest

import (
	"fmt"
	"log"
	"sync"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/redis"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/hasher"
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
			startIdx:         params.StartIdx,
			frameCount:       params.FrameCount,
			jobSize:          configObj.JobSize,
			workers:          make(map[int]process.Process),
			workersListMutex: &sync.Mutex{},
			Cache:            cacheInstance,
			CachePrefix:      params.ExecutionGroupID,
		}

		manager.getQueueNames(configObj.Queues)

		ingestionManagerInstance = &manager
	})

	return ingestionManagerInstance
}

// Start Starts the ingestion manager job scheduling routine
func (manager *IngestionManager) Start() {
	log.Println(logPrefix, "Starting Ingestion Manager")

	log.Println(logPrefix, fmt.Sprintf("Inserting jobs in %s", manager.Queues.Todo))

	jobCount := manager.populateJobsPool()

	log.Println(logPrefix, fmt.Sprintf("Successfully inserted %d jobs in %s", jobCount, manager.Queues.Todo))
}

// Shutdown A function to shutdown the ingestion manager
func (manager *IngestionManager) Shutdown() {
	log.Println(logPrefix, "Ingestion Manager cleaning cache")

	//TODO: Make sure to only do this on success
	manager.flushCache()

	log.Println(logPrefix, "Ingestion Manager shutdown successfully")
}

func (manager *IngestionManager) getQueueNames(queues config.Queue) {
	manager.Queues = Queue{
		Todo:       fmt.Sprintf("%s:%s", manager.CachePrefix, queues.Todo),
		InProgress: fmt.Sprintf("%s:%s", manager.CachePrefix, queues.InProgress),
		Done:       fmt.Sprintf("%s:%s", manager.CachePrefix, queues.Done),
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

		jobToken := hasher.MD5Hash(encodedJob)
		jobs = append(jobs, jobToken)
		jobTokens[jobToken] = encodedJob
	}

	err := manager.insertJobsInQueue(manager.Queues.Todo, jobs...)
	errors.HandleError(err, fmt.Sprintf("%s Unable to insert todo jobs on start up", logPrefix), true)

	err = manager.insertJobTokens(jobTokens)
	errors.HandleError(err, fmt.Sprintf("%s Unable to insert job keys on start up", logPrefix), true)

	return len(jobs)
}
