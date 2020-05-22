package process

import (
	"encoding/json"
)

// NewProcess A function to obtain a process instance
func NewProcess(id int) Process {
	return Process{
		ID:        id,
		Trackable: false,
	}
}

// ParseUtilization A function to parse the process util stats recieved in healthchecks
func ParseUtilization(healthCheck string) (Utilization, error) {
	var util Utilization

	err := json.Unmarshal([]byte(healthCheck), &util)
	if err != nil {
		return Utilization{}, err
	}

	return util, err
}
