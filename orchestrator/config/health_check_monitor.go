package config

// HealthCheckMonitorConfig Houses the configurations of the healthcheck monitor
type HealthCheckMonitorConfig struct {
	IP                       string `yaml:"ip"`                         //IP of the monitor
	Port                     string `yaml:"port"`                       //Port on which it receives healthchecks
	LivenessProbe            int    `yaml:"liveness_probe"`             //A duration (in secs) after which a process is considered dead
	ReadinessProbe           int    `yaml:"readiness_probe"`            //A duration (in secs) to wait for first ping
	HealthCheckInterval      int    `yaml:"health_check_interval"`      //The frequency (in milliseconds) of polling for health checks
	LivenessTrackingInterval int    `yaml:"liveness_tracking_interval"` //The frequency (in milliseconds) for checking dead processes
}

// HealthCheckMonitorConfig A function to return the healthcheck monitor config
func (manager *ConfigurationManager) HealthCheckMonitorConfig(filename string) HealthCheckMonitorConfig {
	var configObj HealthCheckMonitorConfig
	filePath := manager.getFilePath(filename)

	manager.retrieveConfig(&configObj, filePath)

	return configObj
}
