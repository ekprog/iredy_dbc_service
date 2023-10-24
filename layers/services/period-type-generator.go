package services

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/tools"
	"sort"
	"time"
)

type PeriodTypeCallback func(time.Time)

type PeriodTypeProcessor struct {
	log core.Logger
}

func NewPeriodTypeProcessor(log core.Logger) *PeriodTypeProcessor {
	return &PeriodTypeProcessor{
		log: log,
	}
}

// Является ли date одним из точек периода periodType?
func (s *PeriodTypeProcessor) IsMatch(date time.Time, period domain.GenerationPeriod) (bool, error) {

	dateTomorrow := tools.RoundDateTimeToDay(date.UTC().Add(24 * time.Hour))

	nearest, err := s.StepBack(dateTomorrow, period)
	if err != nil {
		return false, errors.Wrap(err, "IsMatch")
	}

	return date.Equal(nearest), nil
}

// Просчитывает время на 1 шаг
// fromDate - с какого дня начинать просчет.
// Текущий день не учитывается!
func (s *PeriodTypeProcessor) StepBack(fromDate time.Time, periodType domain.GenerationPeriod) (time.Time, error) {

	fromDate = tools.RoundDateTimeToDay(fromDate.UTC())

	// Every day
	if periodType.Type == domain.PeriodTypeEveryDay {
		return fromDate.Add(time.Duration(-24) * time.Hour), nil
	}

	// Week day
	if periodType.Type == domain.PeriodTypeWeekDays {
		iDate := fromDate

		// Максимум 7 итераций (если не нашли, то ошибка в данных периода)
		for i := 0; i < 7; i++ {
			iDate = iDate.Add(time.Duration(-24) * time.Hour)
			weekDay := iDate.Weekday()
			for _, item := range periodType.Data {
				if item == int(weekDay) {
					return iDate, nil
				}
			}
		}

		return time.Time{}, errors.New("Incorrect period type Data")
	}

	// Month dates
	if periodType.Type == domain.PeriodTypeMonthDates {
		iDate := fromDate

		// Максимум 31 итераций (если не нашли, то ошибка в данных периода)
		for i := 1; i <= 31; i++ {
			iDate = iDate.Add(time.Duration(-24) * time.Hour)
			monthDay := iDate.Day()
			for _, item := range periodType.Data {
				if item == monthDay {
					return iDate, nil
				}
			}
		}

		return time.Time{}, errors.New("Incorrect month period type Data")
	}

	return time.Time{}, errors.New("Incorrect period type " + periodType.Type)
}

func (s *PeriodTypeProcessor) StepBackN(fromDate time.Time, period domain.GenerationPeriod, n int) (time.Time, error) {
	var err error
	for i := 0; i < n; i++ {
		fromDate, err = s.StepBack(fromDate, period)
		if err != nil {
			return time.Time{}, err
		}
	}
	return fromDate, nil
}

// Итерируется на step шагов назад и вызывает callback с просчитанным временем (сравнение с обрезкой по дню)
func (s *PeriodTypeProcessor) StepBackwardForEach(fromDate time.Time, period domain.GenerationPeriod, step uint, fn PeriodTypeCallback) error {
	var err error
	currentTime := fromDate.UTC()

	for i := uint(0); i < step; i++ {
		// Making one step forward
		currentTime, err = s.StepBack(currentTime, period)
		if err != nil {
			return err
		}

		// Here we have current step. Let's call fn
		fn(currentTime)
	}

	return nil
}

// Возвращает массив последних n итераций, начиная с fromDate (текущий день будет учитываться)
func (s *PeriodTypeProcessor) BackwardList(fromDate time.Time, period domain.GenerationPeriod, n uint) ([]time.Time, error) {

	// Чтобы текущий день тоже учитывался (если он входит в период)
	fromDate = fromDate.UTC().Add(24 * time.Hour)

	var list []time.Time
	err := s.StepBackwardForEach(fromDate, period, n, func(t time.Time) {
		list = append(list, t)
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Окно с пропущенными датами между aDate и bDate
func (s *PeriodTypeProcessor) AbsentWindow(aDate, bDate time.Time, period domain.GenerationPeriod) ([]time.Time, error) {
	aDate = tools.RoundDateTimeToDay(aDate.UTC())
	bDate = tools.RoundDateTimeToDay(bDate.UTC())

	var list []time.Time
	var err error

	for aDate.Before(bDate) {
		bDate, err = s.StepBack(bDate, period)
		if err != nil {
			return nil, errors.Wrap(err, "StepBack")
		}
		if aDate.Before(bDate) {
			list = append(list, bDate)
		}
	}

	// ВАЖНО: возвращаем отсортированный массив дат
	sort.Slice(list, func(i, j int) bool {
		return list[i].Before(list[j])
	})

	return list, nil
}
