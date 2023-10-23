package job

import (
	"github.com/spf13/viper"
	"go.uber.org/dig"
	"reflect"
	"runtime"
)

func NewJob(job interface{}, scheduleUTC string) {
	newJob(job, scheduleUTC, false)
}

func NewJobWithImmediately(job interface{}, scheduleUTC string) {
	newJob(job, scheduleUTC, true)
}

func newJob(job interface{}, scheduleUTC string, immediately bool) {

	name := runtime.FuncForPC(reflect.ValueOf(job).Pointer()).Name()

	scope := di.Scope(name)
	err := scope.Provide(job, dig.As(new(Job)))
	if err != nil {
		log.FatalWrap(err, "cannot initialize new job")
	}

	err = scope.Invoke(func(j Job) {

		runnable := func() error {
			log.Debug("Started %s", name)
			err := j.Run()
			if err != nil {
				return err
			}
			log.Debug("Finished %s", name)
			return nil
		}

		_, err := s.Cron(scheduleUTC).Do(runnable)
		if err != nil {
			log.Fatal("cannot DO cron %s: %s", name, err.Error())
		}

		if immediately {
			immediatelyJobs[name] = runnable
		}
	})
	if err != nil {
		log.FatalWrap(err, "cannot create new job")
	}
	log.Info("New job was successfully registered - %s", name)
}

func Start() error {

	enabled := viper.GetBool("jobs.enabled")
	if !enabled {
		return nil
	}

	// Run immediately
	for name, runnable := range immediatelyJobs {
		err := runnable()
		if err != nil {
			log.Fatal("cannot RUN JOB immediately %s: %s", name, err.Error())
		}
	}

	// Run schedule
	s.StartAsync()
	return nil
}
