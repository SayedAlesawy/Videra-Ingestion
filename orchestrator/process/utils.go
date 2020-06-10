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
		groupArgs := []string{group.Script}
		groupArgs = append(groupArgs, getProcessArgs(group.ArgPrefix, group.ArgAssign, group.Args.Group, true)...)
		groupInstance := newProcessGroup(group.Name, group.Replicas, group.Command, groupArgs)

		for i := 0; i < group.Replicas; i++ {
			processArgs := []string{}
			if i < len(group.Args.Process) {
				processArgs = getProcessArgs(group.ArgPrefix, group.ArgAssign, group.Args.Process[i], false)
			}

			stagedProcessesList = append(stagedProcessesList, newProcess(groupInstance, processArgs))
		}
	}

	return stagedProcessesList
}

// getProcessArgs A function to construct process arguments from config
func getProcessArgs(argsPrefix string, argsAssign string, args []config.ProcessArg, fallback bool) []string {
	processArgs := []string{}

	for _, arg := range args {
		value := arg.Value
		if value == "" && fallback {
			value = params.OrchestratorParamsInstance().ArgsMap[arg.Flag].(string)
		}

		if value == "" {
			continue
		}

		processArgs = append(processArgs, fmt.Sprintf("%s%s%s%s", argsPrefix, arg.Flag, argsAssign, value))
	}

	return processArgs
}
