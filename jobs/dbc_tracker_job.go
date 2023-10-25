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

	challengesRepo domain.DBCChallengesRepository
	tracksRepo     domain.DBCTrackRepository
	usersRepo      domain.UsersRepository
}

func NewDBCTrackerJob(log core.Logger,
	trxManager *manager.Manager,
	pProc *services.PeriodTypeProcessor,
	challengesRepo domain.DBCChallengesRepository,
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
	//date, _ := time.Parse("2006-01-02", "2023-10-13")
	//
	//isDone, err := job.dbcProc.MakeTrack(ctx, challengeId, date, false)
	//if err != nil {
	//	return err
	//}
	//if !isDone {
	//	log.Fatal("isDone == false")
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

		for _, item := range items {
			if item.IsAutoTrack {
				continue
			}
			err := job.dbcProc.ProcessChallengeTracks(ctx, item)
			if err != nil {
				return errors.Wrap(err, "ProcessChallengeTracks")
			}
		}
	}

	return nil
}

//
//func (job *DBCTrackerJob) handleItem(item *domain.DBCChallenge) error {
//	// Рассчитываем дату, до которой нужно обработать челленджи
//	// Учитываем: -2 чтобы попасть на стык
//	timeSince, err := job.pProc.Step(item.CreatedAt, time.Now(), domain.PeriodTypeEveryDay, -3+1)
//	if err != nil {
//		return err
//	}
//
//	// Ищем все неучтенные треки у данного челленджа
//	// Все, что ДО timeSince - подлежит обработке
//	needProcessed, err := job.tracksRepo.FetchNotProcessed(item.Id, timeSince)
//	if err != nil {
//		return err
//	}
//
//	transferScore := 0
//
//	for _, track := range needProcessed {
//		// Сколько очков дает данный трек?
//		// ToDo: пока что всегда 1
//		trackPointWeight := 1
//
//		if track.Done {
//			transferScore += trackPointWeight
//		}
//	}
//
//	tracksIds := lo.Map(needProcessed, func(item *domain.DBCTrack, index int) int64 {
//		return item.Id
//	})
//
//	ctx := context.Background()
//	err = job.trxManager.Do(ctx, func(ctx context.Context) error {
//		// Set tracks as processed
//		err := job.tracksRepo.SetProcessed(ctx, tracksIds)
//		if err != nil {
//			return err
//		}
//
//		// Updating user`s scores
//		err = job.usersRepo.TransferDailyScores(ctx, int64(item.UserId), transferScore)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
