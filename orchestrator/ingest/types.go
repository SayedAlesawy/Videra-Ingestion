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
	Queues            Queue                   //Houses the queues used by the ingestion manager
	Cache             *redis.Client           //Used by the manager to access a persistent caching layer
	CachePrefix       string                  //Prefix for cache keys used for scoping
	JobCount          int                     //Total number of executed jobs
	CheckDoneInterval time.Duration           //Frequency of checking if all jobs are done
}

// Queue Defines a queue used in the ingestion manager
type Queue struct {
	Todo       string //Queue for to be done jobs
	InProgress string //Queue for currently executing jobs
	Done       string //Queue for already executed jobs
}

// ingestionJob Represents the ingestion job params
type ingestionJob struct {
	Jid         int64 `json:"jid"`          //Unique id for job
	StartIdx    int64 `json:"start_idx"`    //Local index to start indexing from within the job range
	FramesCount int64 `json:"frames_count"` //Local frame count to ingest starting from the local startIdx
}
