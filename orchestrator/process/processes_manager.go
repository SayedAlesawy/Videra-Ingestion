package process

import (
	"fmt"
	"log"
	"sync"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Process-Manager]"

// configManagerOnce Used to garauntee thread safety for singleton instances
var processManagerOnce sync.Once

// monitorInstance A singleton instance of the config manager object
var processesManagerInstance *ProcessesManager

// ProcessesManagerInstance A function to return a processes manager instance
func ProcessesManagerInstance() *ProcessesManager {
	processManagerOnce.Do(func() {
		configManager := config.ConfigurationManagerInstance("config/config_files")
		configObj := configManager.ProcessManagerConfig("processes_manager.yaml")

		manager := ProcessesManager{
			stagedProcessesList: buildStagedProcessesList(configObj),
			processesList:       make(map[int]process),
		}

		processesManagerInstance = &manager
	})

	return processesManagerInstance
}

// Start Executes the processes in the staged processes list
func (manager *ProcessesManager) Start() []Process {
	log.Println(logPrefix, fmt.Sprintf("Starting %d processes", len(manager.stagedProcessesList)))

	var executionList []Process

	for i := 0; i < len(manager.stagedProcessesList); i++ {
		//Execute the staged process
		cmd, err := manager.stagedProcessesList[i].execute()
		errors.HandleError(err, fmt.Sprintf("%s Error in ProcessesManager.Start()", logPrefix), false)
		if errors.IsError(err) {
			continue
		}

		//Create an exposed process instance used for tracking
		executedProcess := NewProcess(cmd.Process.Pid)
		log.Println(logPrefix, "Start process under group:", manager.stagedProcessesList[i].Group.Name, "with PID:", executedProcess.Pid)

		//Add the exposed process to the execution list
		executionList = append(executionList, executedProcess)

		//Construct the internal process instance
		manager.stagedProcessesList[i].Handle = cmd
		manager.stagedProcessesList[i].Running = true

		//Construct the internal processes list
		manager.processesList[manager.stagedProcessesList[i].Handle.Process.Pid] = manager.stagedProcessesList[i]
	}

	//Remove all running processes from the staging area
	totalProcesses := len(manager.stagedProcessesList)
	manager.filterStalledProcess()

	//If all processes failed to start, then panic
	if totalProcesses == len(manager.stagedProcessesList) {
		errors.HandleError(errors.New("All processes failed to start"), fmt.Sprintf("%s Error in ProcessesManager.Start()", logPrefix), true)
	}

	log.Println(logPrefix, fmt.Sprintf("Started %d/%d processes successfully", totalProcesses-len(manager.stagedProcessesList), totalProcesses))

	return executionList
}

// Shutdown A function to kill all spawned processes on shutdown
func (manager *ProcessesManager) Shutdown() {
	log.Println(logPrefix, "Processing manager is shutting down")
}

// KillProcess A function to kill a process by its pid
func (manager *ProcessesManager) KillProcess(pid int) error {
	process, exists := manager.processesList[pid]
	if !exists {
		return errors.New(fmt.Sprintf("Process with pid: %d is not found", pid))
	}

	//Kill the running process
	err := process.Handle.Process.Kill()
	if errors.IsError(err) {
		return errors.New(fmt.Sprintf("Unable to kill process under group: %s with Pid: %d", process.Group.Name, pid))
	}

	//Mark the process as not running
	process.Running = false

	//Remove it from active processes list
	delete(manager.processesList, pid)

	//Insert it in staged processes list
	manager.stagedProcessesList = append(manager.stagedProcessesList, process)

	return nil
}

// filterStalledProcess Filter the processes that are not running from the staging area
func (manager *ProcessesManager) filterStalledProcess() {
	var stalledProcesses []process

	for _, process := range manager.stagedProcessesList {
		if !process.Running {
			stalledProcesses = append(stalledProcesses, process)
		}
	}

	manager.stagedProcessesList = stalledProcesses
}
