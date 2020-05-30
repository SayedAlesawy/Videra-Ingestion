package ingest

import "encoding/json"

// newIngestionJob Creates a new ingestion job instance
func newIngestionJob(jid int64, startIdx int64, frameCount int64) ingestionJob {
	return ingestionJob{
		jid:        jid,
		startIdx:   startIdx,
		frameCount: frameCount,
	}
}

// encode Encodes the job into bytes
func (job *ingestionJob) encode() ([]byte, error) {
	return json.Marshal(job)
}
