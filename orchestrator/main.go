package main

import (
	"os"
	"os/signal"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/health"
	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/process"
)

func main() {
	//Intercept os signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	//TODO: Logic to spwan processes
	processList := []process.Process{
		process.NewProcess(1),
		process.NewProcess(2),
	}

	//Init a health check monitor
	monitor := health.MonitorInstance(processList)

	//Start monitor
	monitor.Start()

	select {
	case <-signals:
		monitor.Shutdown()
	}
}
