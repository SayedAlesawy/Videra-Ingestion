package config

// ConfigurationManager An interface for all config objects
type ConfigurationManager struct {
	configFilesDir string //Directroy in which to look for config files
}

// HealthCheckMonitorConfig Houses the configurations of the healthcheck monitor
type HealthCheckMonitorConfig struct {
	IP            string `yaml:"ip"`             //IP of the monitor
	Port          string `yaml:"port"`           //Port on which it receives healthchecks
	LivenessProbe int    `yaml:"liveness_probe"` //A duration (in secs) after which a process is considered dead
}
