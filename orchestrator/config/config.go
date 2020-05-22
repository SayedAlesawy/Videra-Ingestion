package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	"gopkg.in/yaml.v2"
)

// logPrefix Used for hierarchical logging
var logPrefix = "[Configuration-Manager]"

// configManagerOnce Used to garauntee thread safety for singleton instances
var configManagerOnce sync.Once

// monitorInstance A singleton instance of the config manager object
var configManagerInstance *ConfigurationManager

// ConfigurationManagerInstance A function to return a configuration manager instance
func ConfigurationManagerInstance(configFilesDir string) *ConfigurationManager {
	configManagerOnce.Do(func() {
		manager := ConfigurationManager{configFilesDir: configFilesDir}

		configManagerInstance = &manager
	})

	return configManagerInstance
}

// HealthCheckMonitorConfig A function to return the healthcheck monitor config
func (manager *ConfigurationManager) HealthCheckMonitorConfig(filename string) HealthCheckMonitorConfig {
	var configObj HealthCheckMonitorConfig
	filePath := fmt.Sprintf("%s%s", os.ExpandEnv(fmt.Sprintf("%s/", manager.configFilesDir)), filename)

	retrieveConfig(&configObj, filePath)

	return configObj
}

// retrieveConfig A function to read a config file
func retrieveConfig(configObj interface{}, filePath string) {
	configFileContent, err := ioutil.ReadFile(filePath)
	errors.HandleError(err, fmt.Sprintf("%s %s\n", logPrefix, fmt.Sprintf("%s %s", "Unable to read config file:", filePath)), true)

	err = yaml.Unmarshal([]byte(configFileContent), configObj)
	errors.HandleError(err, fmt.Sprintf("%s %s\n", logPrefix, fmt.Sprintf("%s %s", "Unable to unmarshal config file:", filePath)), true)
}
