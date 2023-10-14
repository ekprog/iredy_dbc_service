package bootstrap

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"microservice/app"
	"microservice/app/core"
	"microservice/app/job"
	"microservice/app/kafka"
	"os"
	"os/signal"
	"syscall"
)

func Run(rootPath ...string) error {

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	osSignCh := make(chan os.Signal)
	signal.Notify(osSignCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-osSignCh
		cancel()
	}()

	// ENV, etc
	err := app.InitApp(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while init app")
	}

	// Logger
	logger, err := app.InitLogs(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while init logs")
	}

	// Storage
	err = app.InitStorage()
	if err != nil {
		return errors.Wrap(err, "error while init storage")
	}

	// Database
	db, err := app.InitDatabase()
	if err != nil {
		return errors.Wrap(err, "error while init db")
	}

	// Migrations
	err = app.RunMigrations(rootPath...)
	if err != nil {
		return errors.Wrap(err, "error while making migrations")
	}

	// gRPC
	_, _, err = app.InitGRPCServer()
	if err != nil {
		return errors.Wrap(err, "cannot init gRPC")
	}

	// DI
	di := core.GetDI()

	if err = di.Provide(func() *sql.DB {
		return db
	}); err != nil {
		return errors.Wrap(err, "cannot provide db")
	}

	if err = di.Provide(func() core.Logger {
		return logger
	}); err != nil {
		return errors.Wrap(err, "cannot provide logger")
	}

	// CRON
	err = job.Init(logger, di)
	if err != nil {
		return errors.Wrap(err, "cannot init jobs")
	}

	// KAFKA
	err = kafka.InitKafka(logger)
	if err != nil {
		return errors.Wrap(err, "cannot init kafka")
	}

	// CORE
	if err := initDependencies(di); err != nil {
		return errors.Wrap(err, "error while init dependencies")
	}

	//
	//
	// HERE CORE READY FOR WORK...
	//
	//

	// CRON
	if err := initJobs(); err != nil {
		return errors.Wrap(err, "error while init jobs")
	}

	if err := job.Start(); err != nil {
		return errors.Wrap(err, "error while start jobs")
	}

	// Run gRPC and block
	go app.RunGRPCServer()

	// End context
	<-ctx.Done()

	return nil
}
