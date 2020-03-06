package process

import "time"

// Process Represents the spwan process by the orchestrator
type Process struct {
	ID       int       //Unique process ID
	IP       string    //IP of the machine on which the process runs
	Port     string    //Port of the machine on which the process runs
	LastPing time.Time //The last timestamp sent by the process
}

// NewProcess A function to obtain a process instance
func NewProcess(id int, ip string, port string) Process {
	return Process{
		ID:   id,
		IP:   ip,
		Port: port,
	}
}
