package main

import (
	"encoding/json"

	zmq "github.com/pebbe/zmq4"
)

type mockmsg struct {
	Path      string `json:"path"`
	FraeIndex int    `json:"frameIndex"`
	BatchSize int    `json:"batchSize"`
}

func main() {
	worker, _ := zmq.NewSocket(zmq.REQ)
	defer worker.Close()
	worker.Connect("tcp://localhost:5555")

	var msg mockmsg = mockmsg{
		Path:      "./mm.mp4",
		FraeIndex: 60,
		BatchSize: 20,
	}
	msgBytes, err := json.Marshal(msg)
	println(err)
	for {
		println("sending message")
		println(string(msgBytes))

		worker.SendBytes(msgBytes, 0)

		worker.Recv(0)
	}
}
