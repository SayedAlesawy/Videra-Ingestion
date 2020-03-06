package health

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/tcp"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"github.com/pebbe/zmq4"
)

// Monitor Represents the health check monitor object
type Monitor struct {
	topic             string                  //The topic on which the processes in the process list publish messages
	connectionHandler tcp.Connection          //A TCP socket used for listening to topic
	processList       map[int]process.Process //Records when was each process last seen
	processListMutex  *sync.Mutex             //Used to insure thread safety while accessing the process list
	livenessProbe     time.Duration           //The max allowed delay after which a process is considered dead
}

// NewMonitor A function to obtain a monitor instance
func NewMonitor(processes []process.Process, processListMutex *sync.Mutex, topic string, livenessProbe time.Duration) (Monitor, bool) {
	connection, err := tcp.NewConnection(zmq4.SUB, topic)
	processList := buildProcessList(processes)

	return Monitor{
		topic:             topic,
		connectionHandler: connection,
		processList:       processList,
		processListMutex:  processListMutex,
		livenessProbe:     livenessProbe,
	}, err
}

// Start A function to start the monitor's work
func (monitorObj *Monitor) Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	//Listen to health check signals sent by tracked processes in the process list
	go monitorObj.listenToHealthChecks(&wg)

	//Update the liveness status of processes in the process list
	go monitorObj.trackLiveness(&wg)

	//Wait for both threads to finish execution - waits forever
	wg.Wait()
}

// listenToHealthChecks A function to listen to health checks
func (monitorObj *Monitor) listenToHealthChecks(wg *sync.WaitGroup) {
	defer wg.Done()
	defer monitorObj.connectionHandler.Close()

	monitorObj.establishConnections()

	for {
		healthCheck, _ := monitorObj.connectionHandler.Recv(zmq4.DONTWAIT)

		if healthCheck != "" {
			log.Println("Received", healthCheck)

			go monitorObj.updateProcessLastSeen(healthCheck.(string))
		}
	}
}

// trackLiveness A function to track the alive status of processes in the process list
func (monitorObj *Monitor) trackLiveness(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		monitorObj.processListMutex.Lock()

		for id, process := range monitorObj.processList {
			delay := time.Now().Sub(process.LastPing)

			if delay > monitorObj.livenessProbe {
				connection := tcp.BuildConnectionString(process.IP, process.Port)
				monitorObj.connectionHandler.Disconnect(connection)

				delete(monitorObj.processList, id)

				log.Println(fmt.Sprintf("Process with id: %d went offline", id))

				//TODO: Add logic to kill that process and spawn another
			}
		}

		monitorObj.processListMutex.Unlock()
	}
}

// updateProcessLastSeen A function to update the last seen of processes in the process list
func (monitorObj *Monitor) updateProcessLastSeen(healthCheck string) {
	processID, err := strconv.Atoi((strings.Fields(healthCheck))[1])
	if errors.IsError(err) {
		log.Println("Error in health check ping")
	}

	monitorObj.processListMutex.Lock()
	process := monitorObj.processList[processID]
	process.LastPing = time.Now()
	monitorObj.processList[processID] = process
	monitorObj.processListMutex.Unlock()
}

// establishConnections A function to establish connections with the processes in the process list
func (monitorObj *Monitor) establishConnections() {
	var connections []string

	for _, process := range monitorObj.processList {
		connections = append(connections, tcp.BuildConnectionString(process.IP, process.Port))
	}

	monitorObj.connectionHandler.Connect(connections...)
}

// buildProcessList A function to build the process list out of the processes slice
func buildProcessList(processes []process.Process) map[int]process.Process {
	processList := make(map[int]process.Process)

	for _, process := range processes {
		processList[process.ID] = process
	}

	return processList
}
