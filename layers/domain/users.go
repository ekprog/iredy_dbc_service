package domain

import (
	"context"
	"time"
)

type User struct {
	Id         int64
	Score      int64
	ScoreDaily int64

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type UsersRepository interface {
	FetchById(int64) (*User, error)
	Exist(int64) (bool, error)
	InsertIfNotExists(*User) error
	Remove(int64) error
	Update(*User) error
	TransferDailyScores(ctx context.Context, userId int64, scoreInc int) error
}

type UsersUseCase interface {
	Info(int64) (GetUserResponse, error)
	CreateIfNotExists(*User) (CreateUserResponse, error)
	Remove(int64) (RemoveUserResponse, error)
}

type GetUserResponse struct {
	StatusCode string
	User       User
}

type CreateUserResponse struct {
	StatusCode string
	Id         int64
}

type RemoveUserResponse struct {
	StatusCode string
	Id         int64
}
