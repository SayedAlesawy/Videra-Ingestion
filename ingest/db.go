package ingest

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/SayedAlesawy/Videra-Ingestion/orchestrator/utils/errors"
	_ "github.com/go-sql-driver/mysql" // mysql driver
)

func (manager *IngestionManager) connectDB() *sql.DB {
	log.Println(logPrefix, "starting a new database connection")
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
	errors.HandleError(err, fmt.Sprintf("%s Failed to intialize mysql instance | %s", logPrefix, err), false)

	err = db.Ping()
	errors.HandleError(err, fmt.Sprintf("%s Failed to connect to mysql instance | %s", logPrefix, err), false)
	log.Println(logPrefix, "DB connection established")

	return db
}

// incrementalFieldUpdate updates a field by incrementing its current value by the new value
func (manager *IngestionManager) incrementalFieldUpdate(newValue int, fieldName string) {
	db := manager.connectDB()
	defer db.Close()

	videoToken := manager.videoToken

	queryStatment, err := db.Prepare(fmt.Sprintf("SELECT %s FROM files WHERE token = ?", fieldName))
	errors.HandleError(err, fmt.Sprintf("%s Failed to get table record information | %s", logPrefix, err), false)
	defer queryStatment.Close()

	updateStatment, err := db.Prepare(fmt.Sprintf("UPDATE files SET %s=? WHERE token = ?", fieldName))
	errors.HandleError(err, fmt.Sprintf("%s Failed to get table record information | %s", logPrefix, err), false)
	defer updateStatment.Close()

	oldValue := 0
	err = queryStatment.QueryRow(videoToken).Scan(&oldValue)
	errors.HandleError(err, fmt.Sprintf("%s Failed to fetch old value from table | %s", logPrefix, err), false)

	_, err = updateStatment.Exec(oldValue+newValue, videoToken)
	errors.HandleError(err, fmt.Sprintf("%s Failed to update new value | %s", logPrefix, err), false)
}

// updateProgressIndicator: A function to update number of completed jobs in video table
func (manager *IngestionManager) updateProgressIndicator() {
	log.Println(logPrefix, "updating total number of jobs in db to ", manager.nextJidAssignment-1, " for video ", manager.videoToken)

	manager.incrementalFieldUpdate(1, "total_done_count")
}

// updateTotalJobsIndicator: A function to set the total number of jobs in video table
func (manager *IngestionManager) updateTotalJobsIndicator() {
	log.Println(logPrefix, "updating total number of jobs in db to ", manager.jobCount, " for video ", manager.videoToken)
	manager.incrementalFieldUpdate(manager.jobCount, "total_job_count")
}
