package usecase

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/layers/services"
)

type UsersUseCase struct {
	log       core.Logger
	repo      domain.UsersRepository
	trackProc *services.DBCTrackProcessor
}

func NewUsersUseCase(log core.Logger,
	repo domain.UsersRepository,
	trackProc *services.DBCTrackProcessor) *UsersUseCase {
	return &UsersUseCase{
		log:       log,
		repo:      repo,
		trackProc: trackProc,
	}
}

func (ucase *UsersUseCase) Info(ctx context.Context, id int64) (domain.GetUserResponse, error) {
	user, err := ucase.repo.FetchById(id)
	if err != nil {
		return domain.GetUserResponse{}, errors.Wrap(err, "Info")
	}

	if user == nil {
		return domain.GetUserResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	scores, err := ucase.trackProc.CalculateScores(ctx, id)
	if err != nil {
		return domain.GetUserResponse{}, errors.Wrap(err, "CalculateScores")
	}

	user.Score = scores.Score
	user.ScoreDaily = scores.ScoreDaily

	return domain.GetUserResponse{
		StatusCode: domain.Success,
		User:       *user,
	}, nil
}

func (ucase *UsersUseCase) CreateIfNotExists(user *domain.User) (domain.CreateUserResponse, error) {
	err := ucase.repo.InsertIfNotExists(user)
	if err != nil {
		return domain.CreateUserResponse{}, errors.Wrap(err, "cannot insert user")
	}
	return domain.CreateUserResponse{
		StatusCode: domain.Success,
		Id:         user.Id,
	}, nil
}

func (ucase *UsersUseCase) Remove(userId int64) (domain.RemoveUserResponse, error) {
	err := ucase.repo.Remove(userId)
	if err != nil {
		return domain.RemoveUserResponse{}, errors.Wrap(err, "cannot remove user")
	}
	return domain.RemoveUserResponse{
		StatusCode: domain.Success,
	}, nil
}
