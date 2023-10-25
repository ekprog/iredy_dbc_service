package services

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/pkg/errors"
	"math"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/tools"
	"time"
)

const DBC_MAX_STEP_CAN_CHANGE = 3

type DBCProcessor struct {
	log        core.Logger
	trxManager *manager.Manager

	periodProc          *PeriodTypeProcessor
	challengeRepository domain.DBCChallengesRepository
	trackRepository     domain.DBCTrackRepository
}

func NewDBCTrackProcessor(log core.Logger,
	trxManager *manager.Manager,
	trackProcessor *PeriodTypeProcessor,
	challengeRepository domain.DBCChallengesRepository,
	trackRepository domain.DBCTrackRepository) *DBCProcessor {
	return &DBCProcessor{
		log:                 log,
		periodProc:          trackProcessor,
		challengeRepository: challengeRepository,
		trackRepository:     trackRepository,
		trxManager:          trxManager,
	}
}

// Меняет значение трека и всей предыдущей цепочки треков
// (НЕ ПРОВЕРЯЕТ дату на возможность трека со стороны бизнеса)
func (s *DBCProcessor) MakeTrack(ctx context.Context, challengeId int64, date time.Time, value bool) (bool, error) {

	now := tools.RoundDateTimeToDay(time.Now().UTC())
	date = tools.RoundDateTimeToDay(date.UTC())

	// Нельзя трекать будущие даты
	if date.After(now) {
		return false, nil
	}

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
	match, err := s.periodProc.IsMatch(date, period)
	if err != nil {
		return false, errors.Wrap(err, "IsMatch")
	}
	if !match {
		return false, errors.Wrap(err, "Incorrect date for period")
	}

	//
	//
	//
	lastSeries := int64(0) // Здесь будет посчитанная цепочка (last_series) предыдущего трека (после вставки пропусков если они есть)
	lastScore := int64(0)  // Здесь будет посчитанная цепочка (score) предыдущего трека (после вставки пропусков, если они есть)
	var diff int64

	var dateSince time.Time

	// Находит дату, от которой нужно начинать перерассчет
	firstTrackBefore, err := s.trackRepository.GetLastForChallengeBefore(ctx, challengeId, date)
	if err != nil {
		return false, errors.Wrap(err, "GetLastForChallengeBefore")
	}

	// Мы первый в БД - начинаем с себя
	if firstTrackBefore == nil {
		lastSeries = 0
		lastScore = 0
		dateSince = date.Add(-24 * time.Hour) // Включая текущую дату
	} else {
		lastSeries = firstTrackBefore.LastSeries
		lastScore = firstTrackBefore.Score
		dateSince = firstTrackBefore.Date
	}

	// Получаем окно дат, которые нужно перерассчитать (массив дат будет отсортированный)

	absentDates, err := s.periodProc.AbsentWindow(dateSince, now.Add(24*time.Hour), period)
	if err != nil {
		return false, errors.Wrap(err, "AbsentWindow")
	}

	// Получаем значения треков в данном диапазоне
	tracks, err := s.trackRepository.FetchForChallengeByDates(challenge.Id, absentDates)
	if err != nil {
		return false, errors.Wrap(err, "FetchForChallengeByDates")
	}

	for _, track := range tracks {
		currentValue := track.Done
		if track.Date.Equal(date) {
			currentValue = value
		}

		// Рассчитываем score
		lastScore, lastSeries, diff = s.nextTrackPoints(lastScore, lastSeries, currentValue)

		track.UserId = challenge.UserId
		track.ChallengeId = challenge.Id
		track.LastSeries = lastSeries
		track.Done = currentValue
		track.Score = lastScore
		track.ScoreDaily = diff
	}

	err = s.trackRepository.InsertOrUpdateBulk(ctx, tracks)
	if err != nil {
		return false, errors.Wrap(err, "InsertOrUpdateBulk")
	}

	return true, nil
}

func (s *DBCProcessor) CalculateScores(ctx context.Context, userId int64) (domain.ScorePoints, error) {
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
		dateProcessed, err := s.periodProc.StepBackN(nowDate, period, DBC_MAX_STEP_CAN_CHANGE)
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

// Обрабатывает все треки для учета User.Score и Challenge.LastSeries
func (s *DBCProcessor) ProcessChallengeTracks(ctx context.Context, challenge *domain.DBCChallenge) error {

	return nil

	//period := domain.GenerationPeriod{Type: domain.PeriodTypeEveryDay}
	//
	//date, err := s.getSeparatorDate(period)
	//if err != nil {
	//	return errors.Wrap(err, "getSeparatorDate")
	//}

}

//
// HELPERS
//

// Получение последней даты, которую можно менять (3ий шаг назад)
func (s *DBCProcessor) getSeparatorDate(period domain.GenerationPeriod) (time.Time, error) {
	date := tools.RoundDateTimeToDay(time.Now().UTC().Add(24 * time.Hour))

	backDate, err := s.periodProc.StepBackN(date, period, DBC_MAX_STEP_CAN_CHANGE)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "StepBackN")
	}

	return backDate, nil
}

func (s *DBCProcessor) nextTrackPoints(lastScore int64, lastSeries int64, currentValue bool) (int64, int64, int64) {

	if !currentValue {
		x := int64(math.Floor(float64(lastScore) * 0.2))
		return x, 0, x - lastScore
	} else {
		return lastScore + 1, lastSeries + 1, 1
	}
}
