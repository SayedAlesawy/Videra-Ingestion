package config

// IngestionManagerConfig Houses the configurations of the ingestion manager
type IngestionManagerConfig struct {
	JobSize int64 `yaml:"job_size"` //Frame count per job
	Queues  Queue `yaml:"queues"`   //Defines queues names used in the queueing system
}

// Queue Defines a queue used in the ingestion manager
type Queue struct {
	Todo       string `yaml:"todo"`        //Queue for to be done jobs
	InProgress string `yaml:"in_progress"` //Queue for currently executing jobs
	Done       string `yaml:"done"`        //Queue for already executed jobs
}

// IngestionManagerConfig A function to return the ingestion manager config
func (manager *ConfigurationManager) IngestionManagerConfig(filename string) IngestionManagerConfig {
	var configObj IngestionManagerConfig
	filePath := manager.getFilePath(filename)

	manager.retrieveConfig(&configObj, filePath)

	return configObj
}
