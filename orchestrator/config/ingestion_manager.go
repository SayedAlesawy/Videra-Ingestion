package config

// IngestionManagerConfig Houses the configurations of the ingestion manager
type IngestionManagerConfig struct {
	WorkerPoolIP            string `yaml:"worker_pool_ip"`           //IP of the worker pool
	WorkerPoolPort          string `yaml:"worker_pool_port"`         //Port of the worker pool on which jobs are sent
	WorkersScanningInterval int    `yaml:"worker_scanning_interval"` //The frequency at which the manager checks for non-busy workers
	JobSize                 int64  `yaml:"job_size"`                 //Frame count per job
	JobSendTimeout          int    `yaml:"job_send_timeout"`         //Timeout for sending job to worker pool
}

// IngestionManagerConfig A function to return the ingestion manager config
func (manager *ConfigurationManager) IngestionManagerConfig(filename string) IngestionManagerConfig {
	var configObj IngestionManagerConfig
	filePath := manager.getFilePath(filename)

	manager.retrieveConfig(&configObj, filePath)

	return configObj
}
