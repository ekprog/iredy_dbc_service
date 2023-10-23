package bootstrap

import (
	"database/sql"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/pkg/errors"
	"microservice/app"
	"microservice/app/core"
	"microservice/app/job"
	"microservice/app/kafka"
)

func Run(rootPath ...string) error {

	// ENV, etc
	ctx, _, err := app.InitApp(rootPath...)
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

	if err = di.Provide(func() *trmsql.CtxGetter {
		return trmsql.DefaultCtxGetter
	}); err != nil {
		return errors.Wrap(err, "cannot provide tx getter")
	}

	trManager := manager.Must(
		trmsql.NewDefaultFactory(db),
		manager.WithCtxManager(trmcontext.DefaultManager),
	)
	if err = di.Provide(func() *manager.Manager {
		return trManager
	}); err != nil {
		return errors.Wrap(err, "cannot provide tx manager")
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
