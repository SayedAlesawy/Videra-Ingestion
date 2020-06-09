package ingest

import (
	"fmt"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
)

// insertTodoJobs A function to insert jobs into the todo queue
func (manager *IngestionManager) insertJobsInQueue(queue string, jobs ...interface{}) error {
	return manager.Cache.LPush(queue, jobs...).Err()
}

// findJobInQueue A function to check if a job exists in a queue or not
func (manager *IngestionManager) findJobInQueue(queue string, jobKey string) (bool, error) {
	jobs, err := manager.Cache.LRange(queue, 0, -1).Result()
	if errors.IsError(err) {
		return false, err
	}

	for _, job := range jobs {
		if job == jobKey {
			return true, nil
		}
	}

	return false, nil
}

// A function to get the current active job of a given pid
func (manager *IngestionManager) getActiveJob(pid int) (string, bool) {
	value, err := manager.Cache.HGet(manager.getActiveJobKey(), fmt.Sprintf("%d", pid)).Result()
	if fmt.Sprintf("%v", err) == "redis: nil" && value == "" {
		return "", false
	}

	return value, true
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
}

// getActiveJobKey A function to get the hash name where active jobs are stored
func (manager *IngestionManager) getActiveJobKey() string {
	return fmt.Sprintf("%s:%s", manager.CachePrefix, "ingestion")
}
