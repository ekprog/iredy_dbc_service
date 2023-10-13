package domain

import "time"

type User struct {
	Id int32 `json:"id"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type UsersRepository interface {
	Exist(int32) (bool, error)
	InsertIfNotExists(*User) error
	Remove(int32) error
}

type UsersInteractor interface {
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
