package config

// HealthCheckMonitorConfig Houses the configurations of the healthcheck monitor
type HealthCheckMonitorConfig struct {
	IP                       string `yaml:"ip"`                         //IP of the monitor
	Port                     string `yaml:"port"`                       //Port on which it receives healthchecks
	LivenessProbe            int    `yaml:"liveness_probe"`             //A duration (in secs) after which a process is considered dead
	HealthCheckInterval      int    `yaml:"health_check_interval"`      //The frequency of polling for health checks
	LivenessTrackingInterval int    `yaml:"liveness_tracking_interval"` //The frequency for checking dead processes
}

// HealthCheckMonitorConfig A function to return the healthcheck monitor config
func (manager *ConfigurationManager) HealthCheckMonitorConfig(filename string) HealthCheckMonitorConfig {
	var configObj HealthCheckMonitorConfig
	filePath := manager.getFilePath(filename)

	manager.retrieveConfig(&configObj, filePath)

	return configObj
}
