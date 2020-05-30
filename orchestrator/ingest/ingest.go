package ingest

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/tcp"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/params"
	"github.com/pebbe/zmq4"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Ingestion-Manager]"

// ingestionManagerOnce Used to garauntee thread safety for singleton instances
var ingestionManagerOnce sync.Once

// ingestionManagerInstance A singleton instance of the ingestion manager object
var ingestionManagerInstance *IngestionManager

// ack Ack sent by worker pool
const ack = "ack"

// IngestionManagerInstance A function to return an ingestion manager instance
func IngestionManagerInstance() *IngestionManager {
	ingestionManagerOnce.Do(func() {
		configManager := config.ConfigurationManagerInstance("config/config_files")
		configObj := configManager.IngestionManagerConfig("ingestion_manager.yaml")
		params := params.OrchestratorParamsInstance()

		connection, err := tcp.NewConnection(zmq4.REQ, "")
		errors.HandleError(err, fmt.Sprintf("%s %s\n", logPrefix, "Unable to establish tcp connection"), true)

		activeRoutines := 1

		manager := IngestionManager{
			workerPoolIP:           configObj.WorkerPoolIP,
			workerPoolPort:         configObj.WorkerPoolPort,
			connectionHandler:      connection,
			startIdx:               params.StartIdx,
			frameCount:             params.FrameCount,
			jobSize:                configObj.JobSize,
			workersScaningInterval: time.Duration(configObj.WorkersScanningInterval) * time.Second,
			jobSendTimeout:         time.Duration(configObj.JobSendTimeout) * time.Second,
			jobsListMutex:          &sync.Mutex{},
			workersListMutex:       &sync.Mutex{},
			activeJobsMutex:        &sync.Mutex{},
			activeRoutines:         activeRoutines,
		}

		manager.populateJobsPool()

		ingestionManagerInstance = &manager
	})

	return ingestionManagerInstance
}

// Start Starts the ingestion manager job scheduling routine
func (manager *IngestionManager) Start() {
	log.Println(logPrefix, "Starting Ingestion Manager")

	manager.wg.Add(manager.activeRoutines)

	//Start jobs assigner
	go manager.assignJobs()
}

// Shutdown A function to shutdown the ingestion manager
func (manager *IngestionManager) Shutdown() {
	log.Println(logPrefix, "Shutting down ingestion manager")

	for i := 0; i < manager.activeRoutines; i++ {
		manager.shutdown <- true
	}

	log.Println(logPrefix, "Waiting for ingestion routines to terminate")
	manager.wg.Wait()

	log.Println(logPrefix, "Health check monitor shutdown successfully")
}

func (manager *IngestionManager) assignJobs() {
	defer manager.wg.Done()

	log.Println(logPrefix, "Jobs Assigner Started")

	manager.establishConnection()

	for range time.Tick(manager.workersScaningInterval) {
		select {
		case <-manager.shutdown:
			log.Println(logPrefix, "Ingestion Manager is shutting down")

			return
		default:
			if len(manager.jobsList) == 0 {
				continue
			}

			manager.jobsListMutex.Lock()
			manager.workersListMutex.Lock()
			manager.activeJobsMutex.Lock()

			//Assign jobs to non-busy workers
			for _, worker := range manager.workers {
				if worker.Utilization.Busy {
					continue
				}

				//Check todo jobs
				for _, job := range manager.jobsList {
					go manager.sendJob(job)
					break
				}
			}

			manager.jobsListMutex.Unlock()
			manager.workersListMutex.Unlock()
			manager.activeJobsMutex.Unlock()
		}
	}
}

func (manager *IngestionManager) sendJob(job ingestionJob) {
	encodedJob, err := job.encode()
	errors.HandleError(err, fmt.Sprintf("%s Unable to encode job, err: %s", logPrefix, err), false)
	if errors.IsError(err) {
		return
	}

	log.Println(logPrefix, fmt.Sprintf("Sending job jid: %d to worker pool", job.jid))

	sendChan := make(chan bool, 1)
	go func() {
		sendErr := manager.connectionHandler.Send(encodedJob, 0)
		acknowledge, recvErr := manager.connectionHandler.RecvString(0)
		success := !errors.IsError(sendErr) && !errors.IsError(recvErr) && (acknowledge == ack)

		if success {
			log.Println(logPrefix, fmt.Sprintf("Sending job jid: %d received by worker pool", job.jid))
		} else {
			log.Println(logPrefix, fmt.Sprintf("Unable to send job jid: %d to worker pool", job.jid))
		}

		sendChan <- true
	}()

	select {
	case <-sendChan:
	case <-time.After(manager.jobSendTimeout):
		log.Println(logPrefix, fmt.Sprintf("Sending job jid: %d to worker pool timed out", job.jid))
	}
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
		jid := i + 1

		if remainder != 0 && i == jobsCount-1 {
			jobSize = remainder
		}

		manager.jobsList[jid] = newIngestionJob(jid, start, jobSize)
	}
}

// establishConnection A function to establish connection with the worker pool
func (manager *IngestionManager) establishConnection() {
	endpoint := tcp.BuildConnectionString(manager.workerPoolIP, manager.workerPoolPort)

	manager.connectionHandler.Connect(endpoint)
}
