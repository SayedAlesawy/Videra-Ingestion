package params

import (
	"flag"
	"fmt"
	"sync"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Params-Manager]"

const (
	defaultStringValue string = ""
	defaultIntValue    int    = -1
	defaultLogIntValue int64  = -1
)

// orchestratorParamsOnce Used to garauntee thread safety for singleton instances
var orchestratorParamsOnce sync.Once

// orchestratorParamsInstance A singleton instance of the orchestrator params object
var orchestratorParamsInstance *OrchestratorParams

// OrchestratorParamsInstance A function to return the orchestrator params singleton instance
func OrchestratorParamsInstance() *OrchestratorParams {
	orchestratorParamsOnce.Do(func() {
		videoPath := flag.String("video-path", "", "path to the video to be processed")
		modelPath := flag.String("model-path", "", "path to the model to be applied")
		configPath := flag.String("model-config-path", "", "path to the model config to be applied")
		startIdx := flag.Int64("start-idx", -1, "starting index from which to process video at path")
		frameCount := flag.Int64("frame-count", -1, "number of frames to process starting from start-idx")

		flag.Parse()

		orchestratorParams := OrchestratorParams{
			VideoPath:       validate("video-path", *videoPath).(string),
			ModelPath:       validate("model-path", *modelPath).(string),
			ModelConfigPath: validate("model-config-path", *configPath).(string),
			StartIdx:        validate("start-idx", *startIdx).(int64),
			FrameCount:      validate("frame-count", *frameCount).(int64),
		}

		orchestratorParamsInstance = &orchestratorParams
	})

	return orchestratorParamsInstance
}

// validate A function to validate required params
func validate(param string, value interface{}) interface{} {
	hasDefaultVal := false

	switch value.(type) {
	case string:
		hasDefaultVal = isDefault(value.(string), defaultStringValue)
	case int:
		hasDefaultVal = isDefault(value.(int), defaultIntValue)
	case int64:
		hasDefaultVal = isDefault(value.(int64), defaultLogIntValue)
	default:
		return false
	}

	if hasDefaultVal {
		panic(fmt.Sprintf("Required param: %s not provided", param))
	}

	return value
}

// isDefault Checks if a value is defaulted
func isDefault(value interface{}, defaultVal interface{}) bool {
	return value == defaultVal
}
