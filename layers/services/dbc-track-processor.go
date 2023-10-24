package services

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/pkg/errors"
	"math"
	"microservice/app/core"
	"microservice/layers/domain"
	"time"
)

type DBCTrackProcessor struct {
	log        core.Logger
	trxManager *manager.Manager

	trackProcessor      *PeriodTypeProcessor
	challengeRepository domain.DBCChallengesRepository
	trackRepository     domain.DBCTrackRepository
}

func NewDBCTrackProcessor(log core.Logger,
	trxManager *manager.Manager,
	trackProcessor *PeriodTypeProcessor,
	challengeRepository domain.DBCChallengesRepository,
	trackRepository domain.DBCTrackRepository) *DBCTrackProcessor {
	return &DBCTrackProcessor{
		log:                 log,
		trackProcessor:      trackProcessor,
		challengeRepository: challengeRepository,
		trackRepository:     trackRepository,
		trxManager:          trxManager,
	}
}

// Меняет значение трека и всей предыдущей цепочки треков
// (НЕ ПРОВЕРЯЕТ дату на возможность трека со стороны бизнеса)
func (s *DBCTrackProcessor) MakeTrack(ctx context.Context, challengeId int64, date time.Time, value bool) (bool, error) {

	// Получаем челлендж
	challenge, err := s.challengeRepository.FetchById(ctx, challengeId)
	if err != nil {
		return false, errors.Wrap(err, "FetchById")
	}
	if challenge == nil {
		return false, nil
	}

	// ToDo: Получаем у челленджа период
	period := domain.GenerationPeriod{Type: domain.PeriodTypeEveryDay}

	// Проверяем, что текущий день является точкой периода и может быть трекнут
	match, err := s.trackProcessor.IsMatch(date, period)
	if err != nil {
		return false, errors.Wrap(err, "IsMatch")
	}
	if !match {
		return false, errors.Wrap(err, "Incorrect date for period")
	}

	//
	//
	//
	lastSeries := int64(0) // Здесь будет посчитанная цепочка (last_series) предыдущего трека (после вставки пропуском, если они есть)
	lastScore := int64(0)  // Здесь будет посчитанная цепочка (score) предыдущего трека (после вставки пропуском, если они есть)
	var diff int64

	//
	// Текущий трек может быть сделан без заполнения предыдущих, поэтому важно все пропуски заполнить и
	// перерассчитать цепочку по трекам
	//

	prevTime, err := s.trackProcessor.StepBack(date, period)
	if err != nil {
		return false, errors.Wrap(err, "StepBack")
	}

	prevTrack, err := s.trackRepository.GetByDate(ctx, challengeId, prevTime)
	if err != nil {
		return false, errors.Wrap(err, "GetByDate")
	}

	var absentDates []time.Time
	var absentTracks []*domain.DBCTrack

	// ВАРИАНТ 1: Предыдущий трек есть и мы просто отталкиваемся от него при расчетах
	if prevTrack != nil {
		lastSeries = prevTrack.LastSeries
		lastScore = prevTrack.Score
	}

	// ВАРИАНТ 2: Здесь либо дат ДО вообще нет, либо есть незаполненные пропуски в ДО периоде
	// (обрабатываем обе ситуации)
	if prevTrack == nil {

		firstTrackBefore, err := s.trackRepository.GetLastForChallengeBefore(ctx, challengeId, date)
		if err != nil {
			return false, errors.Wrap(err, "GetLastForChallengeBefore")
		}

		// Значит трекаем первую дату
		if firstTrackBefore == nil {
			lastSeries = 0
			lastScore = 0
		}

		// Тогда заполняем все пропуски и перерассчитываем цепочку
		if firstTrackBefore != nil {
			// Здесь будет отсортированный массив дат, что крайне важно для цепочки
			absentDates, err = s.trackProcessor.AbsentWindow(firstTrackBefore.Date, date, period)
			if err != nil {
				return false, errors.Wrap(err, "AbsentWindow")
			}

			// Рассчитываем цепочку треков на основе окна дат
			lastSeries = firstTrackBefore.LastSeries
			lastScore = firstTrackBefore.Score

			for _, absentDate := range absentDates {

				// Рассчитываем score
				lastScore, lastSeries, diff = s.nextTrackPoints(lastScore, lastSeries, false)

				absentTrack := &domain.DBCTrack{
					UserId:      challenge.UserId,
					ChallengeId: challenge.Id,
					Date:        absentDate,
					Done:        false,      // незаполненный трек будет ложным
					LastSeries:  lastSeries, // после ложного трека серия будет нулевой
					Score:       lastScore,  // после ложного трека пользователь получает 20% от предыдущего
					ScoreDaily:  diff,
				}
				absentTracks = append(absentTracks, absentTrack)
			}
		}
	}

	//
	// lastSeries и lastScore просчитаны для последнего элемента.
	// Все пропуски треков также заполнены.
	// Можно рассчитать новый трек.
	//

	lastSeries, lastScore, diff = s.nextTrackPoints(lastScore, lastSeries, value)
	currentTrack := &domain.DBCTrack{
		UserId:      challenge.UserId,
		ChallengeId: challenge.Id,
		Date:        date,
		Done:        value,
		LastSeries:  lastSeries,
		Score:       lastScore,
		ScoreDaily:  diff,
	}

	//
	// Трек не обязательно является крайним, поэтому важно перерассчитать все треки, которые идут
	// после него
	//
	afterList, err := s.trackRepository.GetAllForChallengeAfter(ctx, challengeId, date)
	if err != nil {
		return false, errors.Wrap(err, "GetAllForChallengeAfter")
	}
	for _, track := range afterList {
		lastScore, lastSeries, diff = s.nextTrackPoints(lastScore, lastSeries, track.Done)
		track.LastSeries = lastSeries
		track.Score = lastScore
		track.ScoreDaily = diff
	}

	// Теперь записываем все критические изменения с учетом транзакций
	// ToDo: можно все отправить одним запросом на вставку/обновление
	err = s.trxManager.Do(ctx, func(ctx context.Context) error {

		// Вставляем в БД недостающую цепочку
		if len(absentTracks) > 0 {
			err = s.trackRepository.InsertNew(ctx, absentTracks)
			if err != nil {
				return errors.Wrap(err, "InsertNew")
			}
		}

		// Вставляем/Обновляем текущий трек
		err = s.trackRepository.InsertOrUpdate(ctx, currentTrack)
		if err != nil {
			return errors.Wrap(err, "InsertOrUpdate")
		}

		// Обновляем все треки ПОСЛЕ
		if len(afterList) > 0 {
			err = s.trackRepository.UpdateSome(ctx, afterList)
			if err != nil {
				return errors.Wrap(err, "UpdateSome")
			}
		}

		return nil
	})
	if err != nil {
		return false, errors.Wrap(err, "TrxManager")
	}

	return true, nil
}

func (s *DBCTrackProcessor) CalculateScores(ctx context.Context, userId int64) (domain.ScorePoints, error) {
	nowDate := time.Now() // ToDo: Location should be like user has (Moscow default)

	// Для каждого челленжда вычисляем scores
	challenges, err := s.challengeRepository.FetchUsersAll(userId)
	if err != nil {
		return domain.ScorePoints{}, errors.Wrap(err, "FetchUsersAll")
	}

	totalScore := int64(0)
	totalDailyScore := int64(0)
	for _, challenge := range challenges {

		// ToDo: real period
		period := domain.GenerationPeriod{Type: domain.PeriodTypeEveryDay}

		// Вычисляем дату на стыке score и dailyScore
		dateProcessed, err := s.trackProcessor.StepBackN(nowDate, period, 3)
		if err != nil {
			return domain.ScorePoints{}, errors.Wrap(err, "StepBackN")
		}

		lastProcessedTrack, err := s.trackRepository.GetByDate(ctx, challenge.Id, dateProcessed)
		if err != nil {
			return domain.ScorePoints{}, errors.Wrap(err, "GetByDate")
		}
		if lastProcessedTrack != nil {
			totalScore += lastProcessedTrack.Score
		}

		//
		dailyTracks, err := s.trackRepository.GetAllForChallengeAfter(ctx, challenge.Id, dateProcessed)
		if err != nil {
			return domain.ScorePoints{}, errors.Wrap(err, "GetAllForChallengeAfter")
		}

		// Последние не заполненные не трогаем
		trueStart := false
		for i := len(dailyTracks) - 1; i >= 0; i-- {
			if dailyTracks[i].Done {
				trueStart = true
			}

			if trueStart {
				totalDailyScore += dailyTracks[i].ScoreDaily
			}
		}
	}

	return domain.ScorePoints{
		Score:      totalScore,
		ScoreDaily: totalDailyScore,
	}, nil
}

func (s *DBCTrackProcessor) nextTrackPoints(lastScore int64, lastSeries int64, currentValue bool) (int64, int64, int64) {

	if !currentValue {
		x := int64(math.Floor(float64(lastScore) * 0.2))
		return x, 0, x - lastScore
	} else {
		return lastScore + 1, lastSeries + 1, 1
	}
}
