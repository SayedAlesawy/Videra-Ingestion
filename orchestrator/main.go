package main

import (
	"os"
	"os/signal"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/health"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/ingest"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
)

func main() {
	//Intercept os signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	//Init a processes manager
	processesManager := process.ProcessesManagerInstance()

	//Execute all processes
	processesList := processesManager.Start()

	//Init a health check monitor
	monitor := health.MonitorInstance(processesList)

	//Register the ingestion manager as a subscriber to monitor data
	monitor.RegisterSubscriber(ingest.IngestionManagerInstance())

	//Start monitor
	monitor.Start()

	select {
	case <-signals:
		processesManager.Shutdown()
		monitor.Shutdown()
	}
}
