package process

import (
	"encoding/json"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// Process Represents the spwan process by the orchestrator
type Process struct {
	ID          int         //Unique process ID
	Trackable   bool        //Indicates if the process is up an running before monitoring begins
	LastPing    time.Time   //The last timestamp sent by the process
	CPUUtil     float32     //Process CPU utilization
	GPUUtil     float32     //Process GPU utilization
	RAMUsage    float32     //Process memroy usage
	Utilization Utilization //Stores utilization stats
}

// Utilization Represents the process util stats received in healthchecks
type Utilization struct {
	PID int     `json:"pid"`
	CPU float32 `json:"cpu"`
	GPU float32 `json:"gpu"`
	RAM float32 `json:"ram"`
}

// NewProcess A function to obtain a process instance
func NewProcess(id int) Process {
	return Process{
		ID:        id,
		Trackable: false,
	}
}

// ParseUtilization A function to parse the process util stats recieved in healthchecks
func ParseUtilization(healthCheck string) (Utilization, bool) {
	var util Utilization

	err := json.Unmarshal([]byte(healthCheck), &util)
	if err != nil {
		return Utilization{}, errors.IsError(err)
	}

	return util, errors.IsError(err)
}
