package ingest

import (
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/tcp"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
)

// IngestionManager Resposible for scheduling ingestion jobs
type IngestionManager struct {
	workerPoolIP           string                  //IP of the worker pool
	connectionHandler      tcp.Connection          //Used to communicate with the worker pool
	startIdx               int64                   //Global index to start indexing from
	frameCount             int64                   //Global frame count to ingest starting from startIdx
	jobSize                int64                   //Frame count per job
	workersScaningInterval time.Duration           //The frequency at which the manager checks for non-busy workers
	jobSendTimeout         time.Duration           //Timeout for sending job to worker pool
	jobsList               map[int64]ingestionJob  //Jobs pool of available jobs
	jobsInFlight           map[int64]bool          //Marks if a job is in-flight or not
	workers                map[int]process.Process //Workers to which the manager assigns jobs
	activeJobs             map[int]ingestionJob    //Dictionary to keep which worker is executing which job
	jobsListMutex          *sync.Mutex             //Used to insure thread safety while accessing the jobs list
	inFlightJobsMutex      *sync.Mutex             //Used to insure thread safety while accessing the in-flight jobs list
	workersListMutex       *sync.Mutex             //Used to insure thread safety while accessing the workers list
	activeJobsMutex        *sync.Mutex             //Used to insure thread safety while accessing the active jobs list
	wg                     sync.WaitGroup          //Used to wait on fired goroutines
	shutdown               chan bool               //Used to handle shutdown signals
	activeRoutines         int                     //Number of active concurrent routines
}

// ingestionJob Represents the ingestion job params
type ingestionJob struct {
	Jid         int64 `json:"jid"`          //Unique id for job
	StartIdx    int64 `json:"start_idx"`    //Local index to start indexing from within the job range
	FramesCount int64 `json:"frames_count"` //Local frame count to ingest starting from the local startIdx
}
