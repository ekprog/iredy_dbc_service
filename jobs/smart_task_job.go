package jobs

import (
	"microservice/app/core"
	"microservice/domain"
)

func SmartTaskJob(log core.Logger, smartUCase domain.SmartTasksInteractor) {
	result, err := smartUCase.GenerateTasks()
	if err != nil {
		log.ErrorWrap(err, "error in SmartTaskJob")
	}
	if result.StatusCode != domain.Success {
		log.ErrorWrap(err, "not successful result in SmartTaskJob")
	}
}
