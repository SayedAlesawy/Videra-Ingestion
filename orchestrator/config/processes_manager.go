package config

// ProcessesManagerConfig Houses the configurations of the processes manager
type ProcessesManagerConfig struct {
	ProcessGroups []ProcessGroupConfig `yaml:"process_groups"` //Represents the process groups the manager is managing
}

// ProcessArg Represents a command line argument to a process
type ProcessArg struct {
	Flag  string `yaml:"flag"`  //Flag used for named args
	Value string `yaml:"value"` //Optional value for the arg. If not provided, it fallsback to gloabl params
}

// ProcessesArgs Represents set of arguments passed to processes
type ProcessesArgs struct {
	Group   []ProcessArg   `yaml:"group"`   //Arguments passed to all processes
	Process [][]ProcessArg `yaml:"process"` //Arguments passed to specific process
}

// ProcessGroupConfig Houses the configurations of a specific process group
type ProcessGroupConfig struct {
	Name      string        `yaml:"name"`       //The name of the process group
	Replicas  int           `yaml:"replicas"`   //Specifies how many process replicas in the process group
	Command   string        `yaml:"command"`    //Speciifes the command used to run the script
	Script    string        `yaml:"script"`     //Specifies the actual script to be run for all processes in the group
	ArgPrefix string        `yaml:"arg_prefix"` //Used to prefix the arg flag, i.e. -flag or --flag
	ArgAssign string        `yaml:"arg_assign"` //Used to assign values to the arg flag, i.e. flag=val or flag val
	Args      ProcessesArgs `yaml:"args"`       //Array of arguments to be passed to the process
}

// ProcessManagerConfig A function to return the processes manager config
func (manager *ConfigurationManager) ProcessManagerConfig(filename string) ProcessesManagerConfig {
	var configObj ProcessesManagerConfig
	filePath := manager.getFilePath(filename)

	manager.retrieveConfig(&configObj, filePath)

	return configObj
}
