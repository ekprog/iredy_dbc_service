package services

import (
	"microservice/app/core"
	"microservice/domain"
	"time"
)

type PeriodMatcherService struct {
	log core.Logger
}

func NewPeriodMatcherService(log core.Logger) *PeriodMatcherService {
	return &PeriodMatcherService{log: log}
}

func (s *PeriodMatcherService) Match(timeX time.Time, timeY time.Time, by domain.GenerationPeriod) (bool, error) {
	return true, nil
}
