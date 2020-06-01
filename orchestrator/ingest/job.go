package ingest

import "encoding/json"

// newIngestionJob Creates a new ingestion job instance
func newIngestionJob(jid int64, startIdx int64, framesCount int64) ingestionJob {
	return ingestionJob{
		Jid:         jid,
		StartIdx:    startIdx,
		FramesCount: framesCount,
	}
}

// encode Encodes the job into bytes
func (job *ingestionJob) encode() ([]byte, error) {
	return json.Marshal(job)
}
