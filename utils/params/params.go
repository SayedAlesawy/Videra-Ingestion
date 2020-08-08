package params

import (
	"flag"
	"fmt"
	"sync"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/file"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Params-Manager]"

// Default values
const (
	defaultStringValue  string = ""
	defaultIntValue     int    = -1
	defaultLongIntValue int64  = -1
)

// Params flag names
const (
	executionGroupIDFlag = "execution-group-id"
	videopathFlag        = "video-path"
	modelPathFlag        = "model-path"
	modelConfigPathFlag  = "model-config-path"
	codePathFlag         = "code-path"
	startIdxFlag         = "start-idx"
	frameCountFlag       = "frame-count"
)

// orchestratorParamsOnce Used to garauntee thread safety for singleton instances
var orchestratorParamsOnce sync.Once

// orchestratorParamsInstance A singleton instance of the orchestrator params object
var orchestratorParamsInstance *OrchestratorParams

// OrchestratorParamsInstance A function to return the orchestrator params singleton instance
func OrchestratorParamsInstance() *OrchestratorParams {
	orchestratorParamsOnce.Do(func() {
		execGroupID := flag.String(executionGroupIDFlag, "", "unique id for the execution group")
		videoPath := flag.String(videopathFlag, "", "path to the video to be processed")
		modelPath := flag.String(modelPathFlag, "", "path to the model to be applied")
		configPath := flag.String(modelConfigPathFlag, "", "path to the model config to be applied")
		codePath := flag.String(codePathFlag, "", "path to the code file provided by user to run model")
		startIdx := flag.Int64(startIdxFlag, -1, "starting index from which to process video at path")
		frameCount := flag.Int64(frameCountFlag, -1, "number of frames to process starting from start-idx")

		flag.Parse()

		orchestratorParams := OrchestratorParams{
			VideoPath:        validatePath(validate(videopathFlag, *videoPath).(string)),
			ModelPath:        validatePath(validate(modelPathFlag, *modelPath).(string)),
			ModelConfigPath:  validatePath(validate(modelConfigPathFlag, *configPath).(string)),
			CodePath:         validate(codePathFlag, *codePath).(string),
			ExecutionGroupID: validate(executionGroupIDFlag, *execGroupID).(string),
			StartIdx:         validate(startIdxFlag, *startIdx).(int64),
			FrameCount:       validate(frameCountFlag, *frameCount).(int64),
		}

		orchestratorParams.mapParams()

		orchestratorParamsInstance = &orchestratorParams
	})

	return orchestratorParamsInstance
}

func (orchParams *OrchestratorParams) mapParams() {
	paramsMap := make(map[string]interface{})

	paramsMap[videopathFlag] = orchParams.VideoPath
	paramsMap[modelPathFlag] = orchParams.ModelPath
	paramsMap[modelConfigPathFlag] = orchParams.ModelConfigPath
	paramsMap[executionGroupIDFlag] = orchParams.ExecutionGroupID
	paramsMap[startIdxFlag] = orchParams.StartIdx
	paramsMap[frameCountFlag] = orchParams.FrameCount
	paramsMap[codePathFlag] = orchParams.CodePath

	orchParams.ArgsMap = paramsMap
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
		hasDefaultVal = isDefault(value.(int64), defaultLongIntValue)
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

// validatePath A function to validate of a path to a file actually exists
func validatePath(path string) string {
	if !file.Exists(path) {
		panic(fmt.Sprintf("File: %s doesn't exist", path))
	}

	return path
}
