package job

import (
	"github.com/go-co-op/gocron"
	"go.uber.org/dig"
	"microservice/app/core"
	"time"
)

var (
	log             core.Logger
	di              *dig.Container
	s               *gocron.Scheduler
	immediatelyJobs map[string]func() error
)

type Job interface {
	Run() error
}

func Init(logger core.Logger, di_ *dig.Container) error {
	log = logger
	di = di_
	s = gocron.NewScheduler(time.UTC)
	s.SingletonModeAll()
	s.WaitForScheduleAll()
	immediatelyJobs = make(map[string]func() error)
	return nil
}
