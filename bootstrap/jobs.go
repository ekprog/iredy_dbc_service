package bootstrap

import (
	"microservice/app/job"
	"microservice/jobs"
)

func initJobs() error {
	job.NewJobWithImmediately(jobs.NewDBCTrackerJob, "0 23 * * *")
	return nil
}
