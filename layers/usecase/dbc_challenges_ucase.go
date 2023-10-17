package usecase

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
	"strings"
)

type ChallengesUseCase struct {
	log            core.Logger
	usersRepo      domain.UsersRepository
	categoryRepo   domain.DBCCategoryRepository
	challengesRepo domain.DBCChallengesRepository
	usersUseCase   domain.UsersUseCase
}

func NewChallengesUseCase(log core.Logger,
	usersRepo domain.UsersRepository,
	projectsRepo domain.DBCCategoryRepository,
	tasksRepo domain.DBCChallengesRepository,
	usersUseCase domain.UsersUseCase) *ChallengesUseCase {
	//
	return &ChallengesUseCase{
		log:            log,
		usersRepo:      usersRepo,
		categoryRepo:   projectsRepo,
		challengesRepo: tasksRepo,
		usersUseCase:   usersUseCase,
	}
}

func (ucase *ChallengesUseCase) All(userId int32) (domain.ChallengesListResponse, error) {
	var items []*domain.DBCChallenge
	var err error

	items, err = ucase.challengesRepo.FetchAll(userId)
	if err != nil {
		return domain.ChallengesListResponse{}, errors.Wrap(err, "cannot fetch dbc-challenges by user id")
	}

	return domain.ChallengesListResponse{
		StatusCode: domain.Success,
		Challenges: items,
	}, nil
}

func (ucase *ChallengesUseCase) Create(form *domain.CreateDBCChallengeForm) (domain.IdResponse, error) {

	//
	_, err := ucase.usersUseCase.CreateIfNotExists(domain.User{Id: form.UserId})
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot insert user before creating task")
	}

	// Is challenge connected to category?
	categoryId := new(int32)
	if form.CategoryName != nil {
		// Finding category
		category, err := ucase.categoryRepo.FetchByName(form.UserId, *form.CategoryName)
		if err != nil {
			return domain.IdResponse{}, errors.Wrap(err, "cannot fetch category before creating task")
		}

		// Creating if not exists
		if category == nil {
			category = &domain.DBCCategory{
				UserId: form.UserId,
				Name:   *form.CategoryName,
			}
			err = ucase.categoryRepo.Insert(category)
			if err != nil {
				return domain.IdResponse{}, errors.Wrap(err, "cannot insert new category before creating task")
			}
		}

		// Set category Id for next step
		*categoryId = category.Id
	}

	// Check if challenge with same name already exists
	form.Name = strings.TrimSpace(form.Name)
	challengeFound, err := ucase.challengesRepo.FetchByName(form.UserId, form.Name)
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot check if challenge exists by name before creating task")
	}
	if challengeFound != nil {
		return domain.IdResponse{
			StatusCode: domain.AlreadyExists,
		}, nil
	}

	// Validation of challenge form
	if form.Name == "" {
		return domain.IdResponse{
			StatusCode: domain.ValidationError,
		}, nil
	}

	// Creating challenge
	challenge := &domain.DBCChallenge{
		UserId:     form.UserId,
		CategoryId: categoryId,
		Name:       form.Name,
		Desc:       form.Desc,
		LastSeries: 0,
	}
	err = ucase.challengesRepo.Insert(challenge)
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot insert new challenge before creating task")
	}

	//
	return domain.IdResponse{
		StatusCode: domain.Success,
		Id:         challenge.Id,
	}, nil
}

func (ucase *ChallengesUseCase) Update(challenge *domain.DBCChallenge) (domain.StatusResponse, error) {

	fetchedChallenge, err := ucase.challengesRepo.FetchById(challenge.Id)
	if err != nil || fetchedChallenge == nil {
		return domain.StatusResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	err = ucase.challengesRepo.Update(challenge)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "cannot update task")
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (ucase *ChallengesUseCase) Remove(userId, taskId int32) (domain.StatusResponse, error) {

	//task, err := ucase.challengesRepo.FetchById(taskId)
	//if err != nil {
	//	return domain.StatusResponse{}, errors.Wrapf(err, "cannot fetch task by id %d", taskId)
	//}
	//
	//if task.UserId != userId {
	//	return domain.StatusResponse{
	//		StatusCode: domain.AccessDenied,
	//	}, nil
	//}
	//
	//err = ucase.challengesRepo.Remove(task.Id)
	//if err != nil {
	//	return domain.StatusResponse{}, errors.Wrap(err, "cannot remove task")
	//}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}
