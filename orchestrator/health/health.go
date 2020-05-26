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
			ip:                       configObj.IP,
			port:                     configObj.Port,
			connectionHandler:        connection,
			processList:              processList,
			processListMutex:         &sync.Mutex{},
			livenessProbe:            time.Duration(configObj.LivenessProbe) * time.Second,
			readinessProbe:           time.Duration(configObj.ReadinessProbe) * time.Second,
			healthCheckInterval:      time.Duration(configObj.HealthCheckInterval) * time.Second,
			livenessTrackingInterval: time.Duration(configObj.LivenessTrackingInterval) * time.Second,
			activeRoutines:           activeRoutines,
			wg:                       sync.WaitGroup{},
			shutdown:                 make(chan bool, activeRoutines),
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

	for range time.Tick(monitorObj.healthCheckInterval) {
		select {
		case <-monitorObj.shutdown:
			log.Println(logPrefix, "Health checker is shutting down")
			updateLastSeenWG.Wait()

			return
		default:
			healthCheck, _ := monitorObj.connectionHandler.RecvString(zmq4.DONTWAIT)

			if healthCheck != "" {
				exists := monitorObj.filterHealthChecks(healthCheck)
				if !exists {
					continue
				}

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

	for range time.Tick(monitorObj.livenessTrackingInterval) {
		select {
		case <-monitorObj.shutdown:
			log.Println(logPrefix, "Liveness tracker is shutting down")
			return
		default:
			monitorObj.processListMutex.Lock()

			for pid, processObj := range monitorObj.processList {
				var reference time.Time
				var threshold time.Duration
				var violation string

				if processObj.Trackable {
					reference = processObj.LastPing
					threshold = monitorObj.livenessProbe
					violation = "violated liveness probe"
				} else {
					reference = processObj.FirstPing
					threshold = monitorObj.readinessProbe
					violation = "violated readiness probe"
				}

				delay := time.Now().Sub(reference)

				if delay > threshold {
					log.Println(logPrefix, fmt.Sprintf("Process with pid: %d %s with delay = %f secs", pid, violation, delay.Seconds()))

					err := process.ProcessesManagerInstance().KillProcess(pid)
					errors.HandleError(err, fmt.Sprintf("%v", err), false)
					if !errors.IsError(err) {
						delete(monitorObj.processList, pid)

						log.Println(logPrefix, fmt.Sprintf("Killed process with pid: %d", pid))
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
	if errors.IsError(err) {
		return
	}

	monitorObj.processListMutex.Lock()
	pid := processUtil.Pid

	process, exists := monitorObj.processList[pid]
	if exists {
		if !process.Trackable {
			process.Trackable = true
		}
		process.Utilization.CPU = processUtil.CPU
		process.Utilization.GPU = processUtil.GPU
		process.Utilization.RAM = processUtil.RAM
		process.LastPing = time.Now()

		monitorObj.processList[pid] = process
	}

	monitorObj.processListMutex.Unlock()
}

func (monitorObj *Monitor) filterHealthChecks(healthCheck string) bool {
	processUtil, err := process.ParseUtilization(healthCheck)
	errors.HandleError(err, fmt.Sprintf("%s%s", logPrefix, "Unable to unmarshal health check ping at filterHealthChecks()"), false)
	if errors.IsError(err) {
		return false
	}

	pid := processUtil.Pid

	monitorObj.processListMutex.Lock()
	_, exists := monitorObj.processList[pid]
	monitorObj.processListMutex.Unlock()

	return exists
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
		processList[process.Pid] = process
	}

	return processList
}
