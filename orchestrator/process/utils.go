package process

import (
	"encoding/json"
	"fmt"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/config"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/params"
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
		groupArgs := getProcessGroupArgs(group.Script, group.ArgPrefix, group.ArgAssign, group.Args)
		groupInstance := newProcessGroup(group.Name, group.Replicas, group.Command, groupArgs)

		for i := 0; i < group.Replicas; i++ {
			stagedProcessesList = append(stagedProcessesList, newProcess(groupInstance))
		}
	}

	return stagedProcessesList
}

// getProcessGroupArgs A function to construct process arguments from config
func getProcessGroupArgs(script string, argsPrefix string, argsAssign string, args []config.ProcessArg) []string {
	groupArgs := []string{script}

	for _, arg := range args {
		value := arg.Value
		if value == "" {
			value = params.OrchestratorParamsInstance().ArgsMap[arg.Flag].(string)
		}

		groupArgs = append(groupArgs, fmt.Sprintf("%s%s%s%s", argsPrefix, arg.Flag, argsAssign, value))
	}

	return groupArgs
}
