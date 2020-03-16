package decoder

import (
	"encoding/json"
	"log"
	"time"

	"gocv.io/x/gocv"
)

// VideoPackage struct to wrap frames
type VideoPackage struct {
	ID      int      `json:"id"`
	Rows    int      `json:"rows"`
	Columns int      `json:"columns"`
	Count   int      `json:"count"`
	Frames  [][]byte `json:"frames"`
}

// VideoDecoder represents decoder for video file
type VideoDecoder struct {
	inputFile string
	video     *gocv.VideoCapture
}

//will be read later from config file, or calculated
const slaveFramesCapacity = 1000

const maxChannelSize = 6

const defaultSleepingSeconds = 3

// NewVideoDecoder creats video decoder instance
func NewVideoDecoder(inputFile string) (VideoDecoder, error) {
	video, err := gocv.VideoCaptureFile(inputFile)
	return VideoDecoder{inputFile: inputFile, video: video}, err
}

// Reset resets opened video in vide decoder to start position
func (decoder *VideoDecoder) Reset() error {
	decoder.Close()
	video, err := gocv.VideoCaptureFile(decoder.inputFile)
	decoder.video = video
	return err
}

// IsOpen checks if a file is opened in video decoder
func (decoder *VideoDecoder) IsOpen() bool {
	if decoder == nil || !decoder.video.IsOpened() {
		return false
	}
	return true
}

// Close closes video in video decoder
func (decoder *VideoDecoder) Close() {
	decoder.video.Close()
}

// DecodeEntire decodes entire video in decoder file, and returns result in channel
func (decoder *VideoDecoder) DecodeEntire(ch chan *string) int {
	log.Printf("Decoding file %v\n", decoder.inputFile)
	if !decoder.IsOpen() {
		return 0
	}

	packageID := 0
	img := gocv.NewMat()
	defer img.Close()

	processedFrames := VideoPackage{Count: 0, Frames: make([][]byte, 0, slaveFramesCapacity)}

	for {

		//if channel has enough message in queue, wait until some messages are dequeued
		for len(ch) == maxChannelSize {
			time.Sleep(defaultSleepingSeconds * time.Second)
		}

		if ok := decoder.video.Read(&img); !ok {

			if processedFrames.Count != 0 {
				sendPackage(&packageID, &processedFrames, ch)
			}

			log.Printf("Finished parsing: %v,\n", decoder.inputFile)
			log.Printf("Total read packages: %v\n", packageID)

			return packageID
		}

		if img.Empty() {
			continue
		}

		processedFrames.Rows = img.Rows()
		processedFrames.Columns = img.Cols()

		processedFrames.Frames = append(processedFrames.Frames, img.ToBytes())
		processedFrames.Count++

		if len(processedFrames.Frames) == slaveFramesCapacity {
			sendPackage(&packageID, &processedFrames, ch)
		}
	}
}

func sendPackage(packageID *int, processedFrames *VideoPackage, ch chan *string) {
	processedFrames.ID = *packageID
	log.Printf("Finished parsing package: %v,\n", *packageID)
	*packageID++
	parsedPackage := parsePackage(processedFrames)
	ch <- parsedPackage
	processedFrames.clear()
}

func parsePackage(frames *VideoPackage) *string {
	parsedBytes, err := json.Marshal(*frames)

	if err != nil {
		log.Println("Error parsing video frames: ", err)
		return nil
	}

	parsedMessage := string(parsedBytes)
	return &parsedMessage
}

func (frames *VideoPackage) clear() {
	frames.Count = 0
	frames.Frames = frames.Frames[:0]
}
