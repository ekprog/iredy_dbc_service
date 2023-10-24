package services

import (
	"context"
	"github.com/pkg/errors"
	"math"
	"microservice/app/core"
	"microservice/layers/domain"
	"time"
)

type DBCTrackProcessor struct {
	log                 core.Logger
	trackProcessor      *PeriodTypeProcessor
	challengeRepository domain.DBCChallengesRepository
	trackRepository     domain.DBCTrackRepository
}

func NewDBCTrackProcessor(log core.Logger,
	trackProcessor *PeriodTypeProcessor,
	challengeRepository domain.DBCChallengesRepository,
	trackRepository domain.DBCTrackRepository) *DBCTrackProcessor {
	return &DBCTrackProcessor{
		log:                 log,
		trackProcessor:      trackProcessor,
		challengeRepository: challengeRepository,
		trackRepository:     trackRepository,
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

			var absentTracks []*domain.DBCTrack
			for _, absentDate := range absentDates {

				// Рассчитываем score
				lastScore, lastSeries = s.nextTrackPoints(lastScore, lastSeries, false)

				absentTrack := &domain.DBCTrack{
					UserId:      challenge.UserId,
					ChallengeId: challenge.Id,
					Date:        absentDate,
					Done:        false,      // незаполненный трек будет ложным
					LastSeries:  lastSeries, // после ложного трека серия будет нулевой
					Score:       lastScore,  // после ложного трека пользователь получает 20% от предыдущего
				}
				absentTracks = append(absentTracks, absentTrack)
			}

			// Вставляем в БД недостающую цепочку
			if len(absentTracks) > 0 {
				err = s.trackRepository.InsertNew(ctx, absentTracks)
				if err != nil {
					return false, errors.Wrap(err, "InsertNew")
				}
			}
		}
	}

	//
	// lastSeries и lastScore просчитаны для последнего элемента.
	// Все пропуски треков также заполнены.
	// Можно вставлять новый трек.
	//

	lastScore, lastSeries = s.nextTrackPoints(lastScore, lastSeries, value)

	err = s.trackRepository.InsertOrUpdate(ctx, &domain.DBCTrack{
		Id:          0,
		UserId:      challenge.UserId,
		ChallengeId: challenge.Id,
		Date:        date,
		Done:        value,
		LastSeries:  lastSeries,
		Score:       lastScore,
	})
	if err != nil {
		return false, errors.Wrap(err, "InsertOrUpdate")
	}

	return true, nil
}

func (s *DBCTrackProcessor) nextTrackPoints(lastScore int64, lastSeries int64, currentValue bool) (int64, int64) {

	if !currentValue {
		return int64(math.Floor(float64(lastScore) * 0.2)), 0
	} else {
		return lastScore + 1, lastSeries + 1
	}
}
