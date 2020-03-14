package main

import (
	"log"
	"sync"
	"time"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/health"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
)

func main() {
	//TODO: Logic to spwan processes
	processList := []process.Process{
		process.NewProcess(1),
		process.NewProcess(2),
	}

	//Mutex for process list thread safety
	var processListMutex sync.Mutex

	//Init a health check monitor
	monitor, err := health.NewMonitor("127.0.0.1", "9092", 2*time.Second, processList, &processListMutex)
	if err {
		log.Panic("Can't init monitor")
	}

	//Start monitor
	monitor.Start()
}
