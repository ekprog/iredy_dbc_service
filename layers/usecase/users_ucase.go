package usecase

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
)

type UsersUseCase struct {
	log  core.Logger
	repo domain.UsersRepository
}

func NewUsersUseCase(log core.Logger, repo domain.UsersRepository) *UsersUseCase {
	return &UsersUseCase{log: log, repo: repo}
}

func (i *UsersUseCase) CreateIfNotExists(user domain.User) (domain.CreateUserResponse, error) {
	err := i.repo.InsertIfNotExists(&user)
	if err != nil {
		return domain.CreateUserResponse{}, errors.Wrap(err, "cannot insert user")
	}
	return domain.CreateUserResponse{
		StatusCode: domain.Success,
		Id:         user.Id,
	}, nil
}

func (i *UsersUseCase) Remove(userId int32) (domain.RemoveUserResponse, error) {
	err := i.repo.Remove(userId)
	if err != nil {
		return domain.RemoveUserResponse{}, errors.Wrap(err, "cannot remove user")
	}
	return domain.RemoveUserResponse{
		StatusCode: domain.Success,
	}, nil
}
