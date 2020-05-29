package ingest

import "encoding/json"

// Job statuses
const (
	todo       = 0
	inProgress = 1
	done       = 2
)

// newIngestionJob Creates a new ingestion job instance
func newIngestionJob(startIdx int64, frameCount int64) ingestionJob {
	return ingestionJob{
		startIdx:   startIdx,
		frameCount: frameCount,
		status:     todo,
	}
}

// encode Encodes the job into bytes
func (job *ingestionJob) encode() ([]byte, error) {
	return json.Marshal(job)
}
