package process

import "time"

// Process Represents the spwan process by the orchestrator
type Process struct {
	ID          int         //Unique process ID
	Trackable   bool        //Indicates if the process is up an running before monitoring begins
	LastPing    time.Time   //The last timestamp sent by the process
	Utilization Utilization //Stores utilization stats
}

// Utilization Represents the process util stats received in healthchecks
type Utilization struct {
	PID int     `json:"pid"`
	CPU float32 `json:"cpu"`
	GPU float32 `json:"gpu"`
	RAM float32 `json:"ram"`
}
