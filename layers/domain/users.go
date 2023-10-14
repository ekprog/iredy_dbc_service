package domain

import "time"

type User struct {
	Id    int32
	Score int32

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt time.Time
}

type UsersRepository interface {
	Exist(int32) (bool, error)
	InsertIfNotExists(*User) error
	Remove(int32) error
}

type UsersUseCase interface {
	CreateIfNotExists(User) (CreateUserResponse, error)
	Remove(int32) (RemoveUserResponse, error)
}

type CreateUserResponse struct {
	StatusCode string
	Id         int32
}

type RemoveUserResponse struct {
	StatusCode string
	Id         int32
}
