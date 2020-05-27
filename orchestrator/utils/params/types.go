package params

// OrchestratorParams Houses the command line params passed to the orchestrator
type OrchestratorParams struct {
	VideoPath       string //Path of the video chunk to be processed
	ModelPath       string //Path of the user CV model to be applied
	ModelConfigPath string //Path of the user CV model config to be applied
	StartIdx        int    //Start index from which processing would start
	FrameCount      int    //Frame count to be processes
}
