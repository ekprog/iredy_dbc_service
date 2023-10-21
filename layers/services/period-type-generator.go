package services

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
	"time"
)

type PeriodTypeCallback func(time.Time)

type PeriodTypeGenerator struct {
	log core.Logger
}

func NewPeriodTypeGenerator(log core.Logger) *PeriodTypeGenerator {
	return &PeriodTypeGenerator{
		log: log,
	}
}

// Просчитывает время на step шагов (может быть и положительное и отрицательное).
// startDate - стартовое время периода (требуется для правильного рассчета).
// nowDate - с какого дня начинать просчет.
// step - шаг вперед или назад, начиная от nowDate.
func (s *PeriodTypeGenerator) Step(startDate, nowDate time.Time, periodType domain.PeriodType, step int) (time.Time, error) {
	if periodType != domain.PeriodTypeEveryDay {
		return time.Time{}, errors.New("incorrect period type")
	}

	// ToDo: 24 - only for PeriodTypeEveryDay!
	steppedTime := nowDate.Add(time.Duration(step*24) * time.Hour)

	return steppedTime, nil
}

// Итерируется на step шагов назад и вызывает callback с просчитанным временем (сравнение с обрезкой по дню)
func (s *PeriodTypeGenerator) StepBackwardForEach(startDate, nowDate time.Time, periodType domain.PeriodType, step uint, fn PeriodTypeCallback) error {
	if periodType != domain.PeriodTypeEveryDay {
		return errors.New("incorrect period type")
	}

	var err error
	currentTime := nowDate

	for i := uint(0); i < step; i++ {
		// Here we have current step. Let's call fn
		fn(currentTime)

		// Making one step forward
		currentTime, err = s.Step(startDate, currentTime, periodType, -1)
		if err != nil {
			return err
		}
	}

	return nil
}

// Возвращает массив последних n итераций, начиная с nowDate
func (s *PeriodTypeGenerator) BackwardList(startDate, nowDate time.Time, periodType domain.PeriodType, n uint) ([]time.Time, error) {
	var list []time.Time
	err := s.StepBackwardForEach(startDate, nowDate, periodType, n, func(t time.Time) {
		list = append(list, t)
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}
