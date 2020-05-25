package process

import (
	"encoding/json"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
)

// ParseUtilization A function to parse the process util stats recieved in healthchecks
func ParseUtilization(healthCheck string) (Utilization, error) {
	var util Utilization

	err := json.Unmarshal([]byte(healthCheck), &util)
	if err != nil {
		return Utilization{}, err
	}

	return util, err
}

// buildStagedProcessesList A function to construct the manager process list
func buildStagedProcessesList(managerConfig config.ProcessesManagerConfig) []process {
	var stagedProcessesList []process

	for _, group := range managerConfig.ProcessGroups {
		groupInstance := newProcessGroup(group.Name, group.Replicas, group.Command)

		for i := 0; i < group.Replicas; i++ {
			stagedProcessesList = append(stagedProcessesList, newProcess(groupInstance))
		}
	}

	return stagedProcessesList
}
