package health

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/drivers/tcp"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"github.com/pebbe/zmq4"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Health-Monitor]"

// monitorOnce Used to garauntee thread safety for singleton instances
var monitorOnce sync.Once

// monitorInstance A singleton instance of the health monitor object
var monitorInstance *Monitor

// Monitor Represents the health check monitor object
type Monitor struct {
	ip                string                  //IP of the monitor
	port              string                  //Port on which the monitor recieves messages from monitored processes
	connectionHandler tcp.Connection          //A TCP socket used for listening to topic
	processList       map[int]process.Process //Records when was each process last seen
	processListMutex  *sync.Mutex             //Used to insure thread safety while accessing the process list
	livenessProbe     time.Duration           //The max allowed delay after which a process is considered dead
	activeRoutines    int                     //The number of active routines the monitor executes
	wg                sync.WaitGroup          //Used to wait on fired goroutines
	shutdown          chan bool               //Used to handle shutdown signals
}

// MonitorInstance A function that returns a health monitor instance
func MonitorInstance(processes []process.Process) *Monitor {
	monitorOnce.Do(func() {
		configManager := config.ConfigurationManagerInstance("config/config_files")
		configObj := configManager.HealthCheckMonitorConfig("health_check_monitor.yaml")

		connection, err := tcp.NewConnection(zmq4.SUB, "")
		errors.HandleError(err, fmt.Sprintf("%s %s\n", logPrefix, "Unable to establish tcp connection"), true)

		processList := buildProcessList(processes)
		activeRoutines := 2

		monitorInstance = &Monitor{
			ip:                configObj.IP,
			port:              configObj.Port,
			connectionHandler: connection,
			processList:       processList,
			processListMutex:  &sync.Mutex{},
			livenessProbe:     time.Duration(configObj.LivenessProbe) * time.Second,
			activeRoutines:    activeRoutines,
			wg:                sync.WaitGroup{},
			shutdown:          make(chan bool, activeRoutines),
		}
	})

	return monitorInstance
}

// IP A function to return the monitor IP
func (monitorObj *Monitor) IP() string {
	return monitorObj.ip
}

// Port A function to return the monitor port
func (monitorObj *Monitor) Port() string {
	return monitorObj.port
}

// ProcessList A function to return the tracked processes list
func (monitorObj *Monitor) ProcessList() map[int]process.Process {
	return monitorObj.processList
}

// Start A function to start the monitor's work
func (monitorObj *Monitor) Start() {
	log.Println(logPrefix, "Starting Monitor")

	monitorObj.wg.Add(monitorObj.activeRoutines)

	//Listen to health check signals sent by tracked processes in the process list
	go monitorObj.listenToHealthChecks()

	//Update the liveness status of processes in the process list
	go monitorObj.trackLiveness()
}

// Shutdown A function to gracefully shutdown the monitor
func (monitorObj *Monitor) Shutdown() {
	log.Println(logPrefix, "Shutting down health check monitor")

	for i := 0; i < monitorObj.activeRoutines; i++ {
		monitorObj.shutdown <- true
	}

	log.Println(logPrefix, "Waiting for monitoring routines to terminate")
	monitorObj.wg.Wait()

	log.Println(logPrefix, "Health check monitor shutdown successfully")
}

// listenToHealthChecks A function to listen to health checks
func (monitorObj *Monitor) listenToHealthChecks() {
	defer monitorObj.wg.Done()
	defer monitorObj.connectionHandler.Close()

	monitorObj.establishConnection()
	log.Println(logPrefix, "Listening for health checks")

	var updateLastSeenWG sync.WaitGroup

	for {
		select {
		case <-monitorObj.shutdown:
			log.Println(logPrefix, "Health checker is shutting down")
			updateLastSeenWG.Wait()

			return
		default:
			healthCheck, _ := monitorObj.connectionHandler.RecvString(zmq4.DONTWAIT)

			if healthCheck != "" {
				log.Println(logPrefix, "Received", healthCheck)

				updateLastSeenWG.Add(1)
				go monitorObj.updateProcessLastSeen(healthCheck, &updateLastSeenWG)
			}
		}
	}
}

// trackLiveness A function to track the alive status of processes in the process list
func (monitorObj *Monitor) trackLiveness() {
	defer monitorObj.wg.Done()

	log.Println(logPrefix, "Tracking processes liveness")

	for {
		select {
		case <-monitorObj.shutdown:
			log.Println(logPrefix, "Liveness tracker is shutting down")
			return
		default:
			monitorObj.processListMutex.Lock()

			for id, process := range monitorObj.processList {
				if process.Trackable {
					delay := time.Now().Sub(process.LastPing)

					if delay > monitorObj.livenessProbe {
						process.Trackable = false
						monitorObj.processList[id] = process

						log.Println(logPrefix, fmt.Sprintf("Process with id: %d went offline", id))

						//TODO: Add logic to kill that process and spawn another
					}
				}
			}

			monitorObj.processListMutex.Unlock()
		}
	}
}

// updateProcessLastSeen A function to update the last seen of processes in the process list
func (monitorObj *Monitor) updateProcessLastSeen(healthCheck string, wg *sync.WaitGroup) {
	defer wg.Done()

	processUtil, err := process.ParseUtilization(healthCheck)
	errors.HandleError(err, fmt.Sprintf("%s%s", logPrefix, "Unable to unmarshal health check ping at updateProcessLastSeen()"), false)

	monitorObj.processListMutex.Lock()
	processID := processUtil.PID

	process := monitorObj.processList[processID]

	if !process.Trackable {
		process.Trackable = true
	}
	process.CPUUtil = processUtil.CPU
	process.GPUUtil = processUtil.GPU
	process.RAMUsage = processUtil.RAM
	process.LastPing = time.Now()

	monitorObj.processList[processID] = process

	monitorObj.processListMutex.Unlock()
}

// establishConnection A function to establish connection with the monitor port
func (monitorObj *Monitor) establishConnection() {
	connection := tcp.BuildConnectionString(monitorObj.ip, monitorObj.port)

	monitorObj.connectionHandler.Bind(connection)
}

// buildProcessList A function to build the process list out of the processes slice
func buildProcessList(processes []process.Process) map[int]process.Process {
	processList := make(map[int]process.Process)

	for _, process := range processes {
		processList[process.ID] = process
	}

	return processList
}
