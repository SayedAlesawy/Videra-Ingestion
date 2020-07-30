package ingest

import (
	"encoding/json"
	"fmt"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// insertTodoJobs A function to insert jobs into the todo queue
func (manager *IngestionManager) insertJobsInQueue(queue string, jobs ...interface{}) error {
	return manager.cache.LPush(queue, jobs...).Err()
}

// insertJobTokens A function to insert job tokens into the the lookup area
func (manager *IngestionManager) insertJobTokens(jobs map[string]string) error {
	return manager.cache.HMSet(manager.getJobTokensHashKey(), jobs).Err()
}

// findJobInQueue A function to check if a job exists in a queue or not
func (manager *IngestionManager) findJobInQueue(queue string, jobToken string) (bool, error) {
	jobs, err := manager.cache.LRange(queue, 0, -1).Result()
	if errors.IsError(err) {
		return false, err
	}

	for _, job := range jobs {
		if job == jobToken {
			return true, nil
		}
	}

	return false, nil
}

// getJobToken A function to get the job details
func (manager *IngestionManager) getJobToken(jid int) (ingestionJob, bool) {
	jobData := ingestionJob{}

	jobToken, err := manager.cache.HGet(manager.getJobTokensHashKey(), fmt.Sprintf("%d", jid)).Result()
	if fmt.Sprintf("%v", err) == "redis: nil" && jobToken == "" {
		return jobData, false
	}

	err = json.Unmarshal([]byte(jobToken), &jobData)
	errors.HandleError(err, "Failed to parse job data", false)
	return jobData, true
}

// getActiveJobToken A function to get the active job token
func (manager *IngestionManager) getActiveJobToken(pid int) (string, bool) {
	jobToken, err := manager.cache.HGet(manager.getActiveJobKey(), fmt.Sprintf("%d", pid)).Result()
	if fmt.Sprintf("%v", err) == "redis: nil" && jobToken == "" {
		return "", false
	}

	return jobToken, true
}

// A function to get the current active job of a given pid
func (manager *IngestionManager) removeActiveJob(pid int) error {
	return manager.cache.HDel(manager.getActiveJobKey(), fmt.Sprintf("%d", pid)).Err()
}

// moveInQueues A function to move an object pointed at by key from src to dst
func (manager *IngestionManager) moveInQueues(key string, src string, dst string) error {
	pipe := manager.cache.TxPipeline()

	pipe.LRem(src, 0, key)
	pipe.LPush(dst, key)

	_, err := pipe.Exec()

	return err
}

// getQueueLength A function to get queue length
func (manager *IngestionManager) getQueueLength(queue string) (int64, error) {
	return manager.cache.LLen(queue).Result()
}

// flushCache A function to flush the queues and the active staging area
func (manager *IngestionManager) flushCache() {
	//Flush queues
	manager.cache.Del(manager.queues.Todo)
	manager.cache.Del(manager.queues.InProgress)
	manager.cache.Del(manager.queues.Done)

	//Flush active jobs area
	manager.cache.Del(manager.getActiveJobKey())

	//Flush job tokens area
	manager.cache.Del(manager.getJobTokensHashKey())
}

// getActiveJobKey A function to get the hash name where active jobs are stored
func (manager *IngestionManager) getActiveJobKey() string {
	return fmt.Sprintf("%s:%s:%s", manager.cachePrefix, "ingestion", "active_jobs")
}

// getJobsKeysHashKey A function to get the hash name where jobs keys are stored
func (manager *IngestionManager) getJobTokensHashKey() string {
	return fmt.Sprintf("%s:%s:%s", manager.cachePrefix, "ingestion", "jobs")
}
