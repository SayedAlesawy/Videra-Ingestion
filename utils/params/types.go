package params

// OrchestratorParams Houses the command line params passed to the orchestrator
type OrchestratorParams struct {
	VideoPath        string                 //Path of the video chunk to be processed
	ModelPath        string                 //Path of the user CV model to be applied
	ModelConfigPath  string                 //Path of the user CV model config to be applied
	CodePath         string                 //Path of the user code provided to run the model
	VideoToken       string                 //unique identifying value of the video
	ExecutionGroupID string                 //Unique ID for the execution group
	StartIdx         int64                  //Start index from which processing would start
	FrameCount       int64                  //Frame count to be processes
	ArgsMap          map[string]interface{} //Houses all params in a key-value format
}
