package decoder

import (
	"encoding/json"
	"log"
	"time"

	"gocv.io/x/gocv"
)

// VideoPackage struct to wrap frames
type VideoPackage struct {
	ID      int      `json:"id"`      //id of package
	Rows    int      `json:"rows"`    //number of rows
	Columns int      `json:"columns"` //number of columns
	Count   int      `json:"count"`   //number of frames
	Frames  [][]byte `json:"frames"`  //decoded frames
}

// VideoDecoder represents decoder for video file
type VideoDecoder struct {
	inputFile       string             //input file name
	video           *gocv.VideoCapture //opencv video capture file
	packageCapacity int                //max number of frames for package
	channelCapacity int                //max channel buffer size, zero or less for infinite buffer
}

const defaultSleepingSeconds = 3 // specifies how long to wait when channel is full

// NewVideoDecoder creats video decoder instance
// channelCapacity specifies max number of frames for package
// packageCapacity specifies max number of frames for package
func NewVideoDecoder(inputFile string, packageCapacity int, channelCapacity int) (VideoDecoder, error) {
	video, err := gocv.VideoCaptureFile(inputFile)
	return VideoDecoder{
		inputFile:       inputFile,
		video:           video,
		packageCapacity: packageCapacity,
		channelCapacity: channelCapacity,
	}, err
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
	if !decoder.video.IsOpened() {
		return false
	}
	return true
}

// Close closes video in video decoder
func (decoder *VideoDecoder) Close() {
	decoder.video.Close()
}

// DecodeEntire decodes entire video, and returns result in channel
// closes channel when done
func (decoder *VideoDecoder) DecodeEntire(ch chan *string) {
	log.Printf("Decoding file %v\n", decoder.inputFile)
	if !decoder.IsOpen() {
		return
	}

	packageID := 0
	img := gocv.NewMat()
	defer img.Close()

	processedFrames := VideoPackage{Count: 0, Frames: make([][]byte, 0, decoder.packageCapacity)}

	decoder.Reset()

	for {
		//if channel has enough message in queue, wait until some messages are dequeued
		for decoder.channelCapacity > 0 && len(ch) == decoder.channelCapacity {
			time.Sleep(defaultSleepingSeconds * time.Second)
		}

		if ok := decoder.video.Read(&img); !ok {

			if processedFrames.Count != 0 {
				sendPackage(&packageID, &processedFrames, ch)
			}

			log.Printf("Finished parsing: %v,\n", decoder.inputFile)
			log.Printf("Total read packages: %v\n", packageID)
			close(ch)
			return
		}

		if img.Empty() {
			continue
		}

		processedFrames.Rows = img.Rows()
		processedFrames.Columns = img.Cols()

		processedFrames.Frames = append(processedFrames.Frames, img.ToBytes())
		processedFrames.Count++

		if len(processedFrames.Frames) == decoder.packageCapacity {
			sendPackage(&packageID, &processedFrames, ch)
		}
	}
}

// DecodePackage decodes certain package from video, and returns result in channel
// closes channel when done
func (decoder *VideoDecoder) DecodePackage(packageID int, ch chan *string) {
	log.Printf("Decoding package %v from file %v\n", packageID, decoder.inputFile)
	if !decoder.IsOpen() {
		return
	}

	img := gocv.NewMat()
	defer img.Close()

	processedFrames := VideoPackage{Count: 0, Frames: make([][]byte, 0, decoder.packageCapacity)}

	decoder.Reset()
	decoder.video.Grab(packageID * decoder.packageCapacity)

	for {

		if ok := decoder.video.Read(&img); !ok {

			if processedFrames.Count != 0 {
				sendPackage(&packageID, &processedFrames, ch)
			}

			log.Printf("Finished parsing package %v from file %v,\n", packageID, decoder.inputFile)
			close(ch)
			return
		}

		if img.Empty() {
			continue
		}

		processedFrames.Rows = img.Rows()
		processedFrames.Columns = img.Cols()

		processedFrames.Frames = append(processedFrames.Frames, img.ToBytes())
		processedFrames.Count++

		if len(processedFrames.Frames) == decoder.packageCapacity {
			sendPackage(&packageID, &processedFrames, ch)
			return
		}
	}
}

// sendPackage is helper function to package frame and send it via channel
func sendPackage(packageID *int, processedFrames *VideoPackage, ch chan *string) {
	processedFrames.ID = *packageID
	log.Printf("Finished parsing package: %v,\n", *packageID)
	*packageID++
	parsedPackage := parsePackage(processedFrames)
	ch <- parsedPackage
	processedFrames.clear()
}

// parsePackage is helper function to parse video package into json format
func parsePackage(frames *VideoPackage) *string {
	parsedBytes, err := json.Marshal(frames)

	if err != nil {
		log.Println("Error parsing video frames: ", err)
		return nil
	}

	parsedMessage := string(parsedBytes)
	return &parsedMessage
}

// clear is helper function to clear video package
func (frames *VideoPackage) clear() {
	frames.Count = 0
	frames.Frames = frames.Frames[:0]
}
