package jobs

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"log"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/layers/services"
	"time"
)

type DBCTrackerJob struct {
	log        core.Logger
	trxManager *manager.Manager

	pProc     *services.PeriodTypeProcessor
	trackProc *services.DBCTrackProcessor

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
	trackProc *services.DBCTrackProcessor) *DBCTrackerJob {
	return &DBCTrackerJob{
		log:            log,
		trxManager:     trxManager,
		pProc:          pProc,
		usersRepo:      usersRepo,
		challengesRepo: challengesRepo,
		tracksRepo:     tracksRepo,
		trackProc:      trackProc,
	}
}

func (job *DBCTrackerJob) Run() error {

	ctx := context.Background()
	challengeId := int64(24)

	date, _ := time.Parse("2006-01-02", "2023-10-25")

	isDone, err := job.trackProc.MakeTrack(ctx, challengeId, date, true)
	if err != nil {
		return err
	}
	if !isDone {
		log.Fatal("isDone == false")
	}

	// Задача: Найти для каждого челленджа все такие неучтенные треки, для которых срок в N последний
	// итераций прошел (пользователю разрешается менять только последние N трек-дня)
	//
	// ВАЖНО 1: данная операция крайне ресурсоемкая, так как обрабатывает каждый челлендж каждого пользователя
	// отдельно (строить таблицу трек-дней), поэтому может занимать значительное время;
	//
	// ВАЖНО 2: важно помнить, что построение данной функции требует безопасность в плане много поточного
	// доступа.
	// 	- При обновлении Score и ScoreDaily не используем жесткое присвоение, а пользуемся +- операциями;

	// Делим на чанки по 1000 и обрабатываем
	//chunkSize := int64(1000)
	//offset := int64(0)
	//for {
	//	items, err := job.challengesRepo.FetchAll(chunkSize, offset)
	//	if err != nil {
	//		return errors.Wrap(err, "DBCTrackerJob")
	//	}
	//	offset += chunkSize
	//
	//	if len(items) == 0 {
	//		break
	//	}
	//
	//	for _, item := range items {
	//
	//		if item.IsAutoTrack {
	//			continue
	//		}
	//
	//		err := job.handleItem(item)
	//		if err != nil {
	//			return errors.Wrap(err, "DBCTrackerJob.handleItem")
	//		}
	//	}
	//}

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
