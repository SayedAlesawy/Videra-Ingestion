package params

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Params-Manager]"

// orchestratorParamsOnce Used to garauntee thread safety for singleton instances
var orchestratorParamsOnce sync.Once

// orchestratorParamsInstance A singleton instance of the orchestrator params object
var orchestratorParamsInstance *OrchestratorParams

// OrchestratorParamsInstance A function to return the orchestrator params singleton instance
func OrchestratorParamsInstance() *OrchestratorParams {
	orchestratorParamsOnce.Do(func() {
		params := os.Args

		orchestratorParams := OrchestratorParams{
			VideoPath:       params[1],
			ModelPath:       params[2],
			ModelConfigPath: params[3],
			StartIdx:        intParam(params[4]),
			FrameCount:      intParam(params[5]),
		}

		orchestratorParamsInstance = &orchestratorParams
	})

	return orchestratorParamsInstance
}

// intParam A function to parse integer params
func intParam(param string) int {
	val, err := strconv.Atoi(param)
	errors.HandleError(err, fmt.Sprintf("%s Unable to covert string param to int", logPrefix), true)

	return val
}
