package domain

import (
	"context"
	"time"
)

// What types of achievements exists?
//
// - After archive some last_series in challenge
// - After archive some points in user profile
type Achievement struct {
	Id          int64
	Title       string
	Desc        string
	TriggerName string

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type AchievementsRepository interface {
	FetchAll(ctx context.Context) ([]*Achievement, error)
}
