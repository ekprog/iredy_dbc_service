package jobs

import (
	"microservice/app/core"
	"microservice/layers/domain"
)

type DBCTrackerJob struct {
	log             core.Logger
	challengesUCase domain.DBCChallengesUseCase
}

func NewDBCTrackerJob(log core.Logger, challengesUCase domain.DBCChallengesUseCase) *DBCTrackerJob {
	return &DBCTrackerJob{
		log:             log,
		challengesUCase: challengesUCase,
	}
}

func (job *DBCTrackerJob) Run() error {
	job.log.Info("Hello")
	return nil
}
