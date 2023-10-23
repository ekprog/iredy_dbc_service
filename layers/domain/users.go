package domain

import (
	"context"
	"time"
)

type User struct {
	Id         int32
	Score      int32
	ScoreDaily int32

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type UsersRepository interface {
	FetchById(int32) (*User, error)
	Exist(int32) (bool, error)
	InsertIfNotExists(*User) error
	Remove(int32) error
	Update(*User) error
	TransferDailyScores(ctx context.Context, userId int64, scoreInc int) error
}

type UsersUseCase interface {
	Info(int32) (GetUserResponse, error)
	CreateIfNotExists(*User) (CreateUserResponse, error)
	Remove(int32) (RemoveUserResponse, error)
}

type GetUserResponse struct {
	StatusCode string
	User       User
}

type CreateUserResponse struct {
	StatusCode string
	Id         int32
}

type RemoveUserResponse struct {
	StatusCode string
	Id         int32
}
