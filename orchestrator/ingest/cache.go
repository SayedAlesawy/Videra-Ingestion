package ingest

import (
	"fmt"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// insertTodoJobs A function to insert jobs into the todo queue
func (manager *IngestionManager) insertJobsInQueue(queue string, jobs ...interface{}) error {
	return manager.Cache.LPush(queue, jobs...).Err()
}

// insertJobTokens A function to insert job tokens into the the lookup area
func (manager *IngestionManager) insertJobTokens(jobs map[string]string) error {
	return manager.Cache.HMSet(manager.getJobTokensHashKey(), jobs).Err()
}

// findJobInQueue A function to check if a job exists in a queue or not
func (manager *IngestionManager) findJobInQueue(queue string, jobToken string) (bool, error) {
	jobs, err := manager.Cache.LRange(queue, 0, -1).Result()
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

// getActiveJobToken A function to get the active job token
func (manager *IngestionManager) getActiveJobToken(pid int) (string, bool) {
	jobToken, err := manager.Cache.HGet(manager.getActiveJobKey(), fmt.Sprintf("%d", pid)).Result()
	if fmt.Sprintf("%v", err) == "redis: nil" && jobToken == "" {
		return "", false
	}

	return jobToken, true
}

// A function to get the current active job of a given pid
func (manager *IngestionManager) removeActiveJob(pid int) error {
	return manager.Cache.HDel(manager.getActiveJobKey(), fmt.Sprintf("%d", pid)).Err()
}

// moveInQueues A function to move an object pointed at by key from src to dst
func (manager *IngestionManager) moveInQueues(key string, src string, dst string) error {
	pipe := manager.Cache.TxPipeline()

	pipe.LRem(src, 0, key)
	pipe.LPush(dst, key)

	_, err := pipe.Exec()

	return err
}

// flushCache A function to flush the queues and the active staging area
func (manager *IngestionManager) flushCache() {
	//Flush queues
	manager.Cache.Del(manager.Queues.Todo)
	manager.Cache.Del(manager.Queues.InProgress)
	manager.Cache.Del(manager.Queues.Done)

	//Flush active jobs area
	manager.Cache.Del(manager.getActiveJobKey())

	//Flush job tokens area
	manager.Cache.Del(manager.getJobTokensHashKey())
}

// getActiveJobKey A function to get the hash name where active jobs are stored
func (manager *IngestionManager) getActiveJobKey() string {
	return fmt.Sprintf("%s:%s:%s", manager.CachePrefix, "ingestion", "active_jobs")
}

// getJobsKeysHashKey A function to get the hash name where jobs keys are stored
func (manager *IngestionManager) getJobTokensHashKey() string {
	return fmt.Sprintf("%s:%s:%s", manager.CachePrefix, "ingestion", "jobs")
}
