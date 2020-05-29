package health

import (
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/tcp"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/pubsub"
)

// Monitor Represents the health check monitor object
type Monitor struct {
	ip                       string                  //IP of the monitor
	port                     string                  //Port on which the monitor recieves messages from monitored processes
	connectionHandler        tcp.Connection          //A TCP socket used for listening to topic
	processList              map[int]process.Process //Records when was each process last seen
	processListMutex         *sync.Mutex             //Used to insure thread safety while accessing the process list
	livenessProbe            time.Duration           //The max allowed delay after which a process is considered dead
	readinessProbe           time.Duration           //The max allowed delay before the process sends its first ping on startup
	healthCheckInterval      time.Duration           //The frequency at which the monitor polls for healthchecks
	livenessTrackingInterval time.Duration           //The frequency at which the monitor checks dead processes
	activeRoutines           int                     //The number of active routines the monitor executes
	wg                       sync.WaitGroup          //Used to wait on fired goroutines
	shutdown                 chan bool               //Used to handle shutdown signals
	subscribers              []pubsub.Subscriber     //Array of subscribers to publish monitor data for them
}
