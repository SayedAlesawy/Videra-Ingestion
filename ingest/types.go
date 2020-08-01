package ingest

import (
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/redis"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
)

// IngestionManager Resposible for scheduling ingestion jobs
type IngestionManager struct {
	startIdx          int64                   //Global index to start indexing from
	frameCount        int64                   //Global frame count to ingest starting from startIdx
	jobSize           int64                   //Frame count per job
	workers           map[int]process.Process //Workers to which the manager assigns jobs
	workersListMutex  *sync.Mutex             //Used to insure thread safety while accessing the workers list
	queues            Queue                   //Houses the queues used by the ingestion manager
	cache             *redis.Client           //Used by the manager to access a persistent caching layer
	cachePrefix       string                  //Prefix for cache keys used for scoping
	jobCount          int                     //Total number of executed jobs
	checkDoneInterval time.Duration           //Frequency of checking if all jobs are done
}

// Queue Defines a queue used in the ingestion manager
type Queue struct {
	Todo       string //Queue for to be done jobs
	InProgress string //Queue for currently executing jobs
	Done       string //Queue for already executed jobs
}

// ingestionJob Represents the ingestion job params
type ingestionJob struct {
	Jid         int64  `json:"jid"`          //Unique id for job
	StartIdx    int64  `json:"start_idx"`    //Local index to start indexing from within the job range
	FramesCount int64  `json:"frames_count"` //Local frame count to ingest starting from the local startIdx
	Action      string `json:"action"`       //type of action to be executed on the data range specified
}

// Job Types
const (
	mergeAction   = "merge"   // merge model ouptut in periods statsifying constraints
	executeAction = "execute" //Execute model on frames
	nullAction    = "null"    // nothing to do here
)

var actionPipeline = map[string]string{
	executeAction: mergeAction,
	mergeAction:   nullAction,
	nullAction:    nullAction,
}