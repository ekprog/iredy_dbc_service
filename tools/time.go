package tools

import (
	"github.com/pkg/errors"
	"time"
)

func RoundDateTimeToDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func RoundDateTimeToMonth(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 0, 0, 0, 0, 0, time.UTC)
}

func IsEqualDateTimeByDay(t1, t2 time.Time) bool {
	return RoundDateTimeToDay(t1).Equal(RoundDateTimeToDay(t2))
}

func ParseISO(dateISO string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, dateISO)

	if err != nil {
		return time.Time{}, errors.Wrap(err, "parseISO")
	}
	return t, nil
}
