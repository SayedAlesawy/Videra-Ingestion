package main

import (
	"encoding/json"
	"log"
	"os"

	"gocv.io/x/gocv"
)

//VideoFrames struct to wrap frames information
type VideoFrames struct {
	Rows    int      `json:"rows"`
	Columns int      `json:"columns"`
	Count   int      `json:"count"`
	Frames  [][]byte `json:"frames"`
}

func main() {
	//Todo: pass input file name
	inputFile := os.Args[1]

	video, err := gocv.VideoCaptureFile(inputFile)
	if err != nil {
		log.Printf("Error opening video capture file: %s\n", inputFile)
		return
	}

	img := gocv.NewMat()
	defer img.Close()

	//will be read later from config file, or calculated
	slaveFramesCapacity := 100
	processedFrames := VideoFrames{Count: 0, Frames: make([][]byte, 0, slaveFramesCapacity)}

	var sentMessage *string
	for {
		if ok := video.Read(&img); !ok {
			log.Printf("Finished parsing: %v\n", inputFile)
			return
		}
		if img.Empty() {
			continue
		}

		processedFrames.Rows = img.Rows()
		processedFrames.Columns = img.Cols()
		processedFrames.Frames = append(processedFrames.Frames, img.ToBytes())
		processedFrames.Count++

		if len(processedFrames.Frames) == slaveFramesCapacity {
			sentMessage = parseMessage(&processedFrames)
			sendMessage(sentMessage)
			processedFrames.clear()
		}
	}
}

func parseMessage(frames *VideoFrames) *string {
	parsedBytes, err := json.Marshal(*frames)

	if err != nil {
		log.Println("Error parsing video frames: ", err)
		return nil
	}

	parsedMessage := string(parsedBytes)
	return &parsedMessage
}

func (frames *VideoFrames) clear() {
	frames.Count = 0
	frames.Frames = frames.Frames[:0]
}

func sendMessage(message *string) {
	//todo
}
