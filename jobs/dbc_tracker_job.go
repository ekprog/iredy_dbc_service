package jobs

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/layers/services"
)

type DBCTrackerJob struct {
	log        core.Logger
	trxManager *manager.Manager

	pProc   *services.PeriodTypeProcessor
	dbcProc *services.DBCProcessor

	challengesRepo domain.DBCUserChallengeRepository
	tracksRepo     domain.DBCTrackRepository
	usersRepo      domain.UsersRepository
}

func NewDBCTrackerJob(log core.Logger,
	trxManager *manager.Manager,
	pProc *services.PeriodTypeProcessor,
	challengesRepo domain.DBCUserChallengeRepository,
	tracksRepo domain.DBCTrackRepository,
	usersRepo domain.UsersRepository,
	trackProc *services.DBCProcessor) *DBCTrackerJob {
	return &DBCTrackerJob{
		log:            log,
		trxManager:     trxManager,
		pProc:          pProc,
		usersRepo:      usersRepo,
		challengesRepo: challengesRepo,
		tracksRepo:     tracksRepo,
		dbcProc:        trackProc,
	}
}

func (job *DBCTrackerJob) Run() error {

	ctx := context.Background()
	//challengeId := int64(24)
	//
	//date, _ := time.Parse("2006-01-02", "2023-10-26")
	//
	//isDone, err := job.dbcProc.MakeTrack(ctx, challengeId, date, false)
	//if err != nil {
	//	return err
	//}
	//if !isDone {
	//	panic("is DONE false")
	//}
	//
	//return nil

	//Делим на чанки по 1000 и обрабатываем
	chunkSize := int64(1000)
	offset := int64(0)
	for {
		items, err := job.challengesRepo.FetchAll(chunkSize, offset)
		if err != nil {
			return errors.Wrap(err, "FetchAll")
		}
		offset += chunkSize

		if len(items) == 0 {
			break
		}

		var errorList []error
		for _, item := range items {
			// ToDo: ошибка для одного пользователя прерывает все?
			// Also check what we should do with error list
			// ETK Stack ?
			if item.IsAutoTrack {
				err := job.dbcProc.ProcessAutoChallengeTracks(ctx, item)
				if err != nil {
					errorList = append(errorList, errors.Wrap(err, "ProcessAutoChallengeTracks"))
					continue
				}
			} else {
				err := job.dbcProc.ProcessChallengeTracks(ctx, item)
				if err != nil {
					errorList = append(errorList, errors.Wrap(err, "ProcessChallengeTracks"))
					continue
				}
			}
		}
	}

	return nil
}
