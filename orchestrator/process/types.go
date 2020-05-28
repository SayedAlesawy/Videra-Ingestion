package process

import (
	"os/exec"
	"time"
)

// ProcessesManager Manages a set of processes (spawn, kill, etc.)
type ProcessesManager struct {
	stagedProcessesList []process       //Array of staged processes to be executed on run
	processesList       map[int]process //Processes managed by the manager instance
}

// Process Represents the spwan process by the orchestrator
type Process struct {
	Trackable   bool        //Indicates if the process is up an running before monitoring begins
	FirstPing   time.Time   //The the timestamp of the first ping by the process
	LastPing    time.Time   //The the timestamp of the last ping by the process
	Utilization Utilization //Stores utilization stats
	Pid         int         //Unique process ID
}

// Utilization Represents the process util stats received in healthchecks
type Utilization struct {
	Pid int     `json:"pid"` //Specifies a certain process
	CPU float32 `json:"cpu"` //CPU utilization
	GPU float32 `json:"gpu"` //GPU utilization
	RAM float32 `json:"ram"` //RAM utilization
}

// process Represents the staged process structure which is more privileged than the exposed type
type process struct {
	Handle  *exec.Cmd    //A handle on the command that executed the process
	Group   processGroup //The process group to which the process belongs
	Running bool         //Indicates if the process is currently running
}

// processGroup Represents a group of similar processes
type processGroup struct {
	Name     string   //The name for the process group
	Replicas int      //The number of replicas
	Command  string   //Path to the process group command
	Script   string   //Path to the process group script
	Args     []string //Array of arguments to be passed to the process
}
