package config

// ProcessesManagerConfig Houses the configurations of the processes manager
type ProcessesManagerConfig struct {
	ProcessGroups []ProcessGroupConfig `yaml:"process_groups"` //Represents the process groups the manager is managing
}

// ProcessGroupConfig Houses the configurations of a specific process group
type ProcessGroupConfig struct {
	Name     string `yaml:"name"`     //The name of the process group
	Replicas int    `yaml:"replicas"` //Specifies how many process replicas in the process group
	Command  string `yaml:"command"`  //Speciifes the path of the entry point for each process in the group
}

// ProcessManagerConfig A function to return the processes manager config
func (manager *ConfigurationManager) ProcessManagerConfig(filename string) ProcessesManagerConfig {
	var configObj ProcessesManagerConfig
	filePath := manager.getFilePath(filename)

	manager.retrieveConfig(&configObj, filePath)

	return configObj
}
