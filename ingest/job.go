package ingest

import (
	"encoding/json"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// newIngestionJob Creates a new ingestion job instance
func newIngestionJob(jid int64, startIdx int64, framesCount int64, action string) ingestionJob {
	return ingestionJob{
		Jid:         jid,
		StartIdx:    startIdx,
		FramesCount: framesCount,
		Action:      action,
	}
}

// encode Encodes the job into bytes
func (job *ingestionJob) encode() (string, error) {
	encodedJob, err := json.Marshal(job)
	if errors.IsError(err) {
		return "", err
	}

	return string(encodedJob), nil
}

// decode Decodes the stringified job
func decode(encodedJob string) (ingestionJob, error) {
	var job ingestionJob

	err := json.Unmarshal([]byte(encodedJob), &job)
	if errors.IsError(err) {
		return ingestionJob{}, err
	}

	return job, nil
}
