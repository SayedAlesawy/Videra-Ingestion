package ingest

import (
	"sync"

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
		params := params.OrchestratorParamsInstance()

		manager := IngestionManager{
			startIdx:   params.StartIdx,
			frameCount: params.FrameCount,
			jobSize:    10,
		}

		manager.populateJobsPool()

		ingestionManagerInstance = &manager
	})

	return ingestionManagerInstance
}

// Start Starts the ingestion manager job scheduling routine
func (manager *IngestionManager) Start() {
	//TODO
}

// Notify A function to notify the ingestion manager with changes in the process list
func (manager *IngestionManager) Notify(notification interface{}) {
	//processesList := notification.(map[int]process.Process)
	//TODO
}

// populateJobsPool Populates the jobs pool of the ingestion manager
func (manager *IngestionManager) populateJobsPool() {
	jobsCount := manager.frameCount / manager.jobSize
	remainder := manager.frameCount % manager.jobSize

	if remainder != 0 {
		jobsCount++
	}

	for i, start := int64(0), manager.startIdx; i < jobsCount; i, start = i+1, start+manager.jobSize {
		jobSize := manager.jobSize

		if remainder != 0 && i == jobsCount-1 {
			jobSize = remainder
		}

		manager.jobs = append(manager.jobs, newIngestionJob(start, jobSize))
	}
}
