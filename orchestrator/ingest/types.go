package ingest

// IngestionManager Resposible for scheduling ingestion jobs
type IngestionManager struct {
	startIdx   int64          //Global index to start indexing from
	frameCount int64          //Global frame count to ingest starting from startIdx
	jobSize    int64          //Frame count per job
	jobs       []ingestionJob //Jobs pool of available jobs
}

// ingestionJob Represents the ingestion job params
type ingestionJob struct {
	startIdx   int64 //Local index to start indexing from within the job range
	frameCount int64 //Local frame count to ingest starting from the local startIdx
	status     int   //Represents the job status, 0: todo, 1: in-progress, 2: done
}
