package process

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// execute A function to execute a staged process
func (processObj *process) execute() (*exec.Cmd, error) {

	args := []string{}
	args = append(args, processObj.Group.Args...)
	args = append(args, processObj.Args...)
	cmd := exec.Command(processObj.Group.Command, args...)

	err := cmd.Start()
	if errors.IsError(err) {
		return nil, errors.New(fmt.Sprintf("Unable to start process under group: %s", processObj.Group.Name))
	}

	return cmd, nil
}

// NewProcess A function to create a new exposed process instance
func NewProcess(pid int) Process {
	return Process{
		Pid:       pid,
		Trackable: false,
		FirstPing: time.Now(),
	}
}

// newProcess A function to create a new internal process instance
func newProcess(group processGroup, args []string) process {
	return process{
		Group:   group,
		Args:    args,
		Running: false,
	}
}

// newGroup A function to create a new process group instance
func newProcessGroup(name string, replicas int, command string, args []string) processGroup {
	return processGroup{
		Name:     name,
		Replicas: replicas,
		Command:  command,
		Args:     args,
	}
}
