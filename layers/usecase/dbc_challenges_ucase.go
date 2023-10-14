package usecase

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
)

type ChallengesUseCase struct {
	log            core.Logger
	usersRepo      domain.UsersRepository
	categoryRepo   domain.DBCCategoryRepository
	challengesRepo domain.DBCChallengesRepository
}

func NewChallengesUseCase(log core.Logger,
	usersRepo domain.UsersRepository,
	projectsRepo domain.DBCCategoryRepository,
	tasksRepo domain.DBCChallengesRepository) *ChallengesUseCase {
	return &ChallengesUseCase{log: log, usersRepo: usersRepo, categoryRepo: projectsRepo, challengesRepo: tasksRepo}
}

func (i *ChallengesUseCase) All(userId int32) (domain.ChallengesListResponse, error) {
	var items []*domain.DBCChallenge
	var err error

	items, err = i.challengesRepo.FetchAll(userId)
	if err != nil {
		return domain.ChallengesListResponse{}, errors.Wrap(err, "cannot fetch dbc-challenges by user id")
	}

	return domain.ChallengesListResponse{
		StatusCode: domain.Success,
		Challenges: items,
	}, nil
}

func (i *ChallengesUseCase) Create(form *domain.CreateDBCChallengeForm) (domain.IdResponse, error) {

	// ToDo:
	// Если нет категории у пользователя  CategoryName - создать ее
	// Проверить, нет ли челленджа с таким же именем
	// Обрезать у челленджа пробелы
	// Далее создать челлендж и привязать к категории (если она есть)

	/// НИЖЕ СТАРАЯ

	//if form.Name == "" {
	//	return domain.IdResponse{
	//		StatusCode: domain.ValidationError,
	//	}, nil
	//}
	//
	//// If user does not exist - create
	//err := i.usersRepo.InsertIfNotExists(&domain.User{
	//	Id: task.UserId,
	//})
	//if err != nil {
	//	return domain.IdResponse{}, errors.Wrap(err, "cannot insert user before creating task")
	//}
	//
	//// If task is a child of project
	//if task.CategoryId != nil {
	//	project, err := i.categoryRepo.FetchById(*task.CategoryId)
	//	if err != nil {
	//		return domain.IdResponse{}, errors.Wrapf(err, "cannot fetch project %d", *task.CategoryId)
	//	}
	//	if project == nil {
	//		return domain.IdResponse{
	//			StatusCode: domain.ProjectNotFound,
	//		}, nil
	//	}
	//	// Is user owner of project?
	//	if project.UserId != task.UserId {
	//		return domain.IdResponse{
	//			StatusCode: domain.AccessDenied,
	//		}, nil
	//	}
	//}
	//
	//err = i.challengesRepo.Insert(task)
	//if err != nil {
	//	return domain.IdResponse{}, errors.Wrap(err, "cannot insert task")
	//}
	//
	//return domain.IdResponse{
	//	StatusCode: domain.Success,
	//	Id:         task.Id,
	//}, nil
	return domain.IdResponse{}, nil
}

func (i *ChallengesUseCase) Update(task *domain.DBCChallenge) (domain.StatusResponse, error) {

	// Поменять имя
	//
	//fetchedTask, err := i.challengesRepo.FetchById(task.Id)
	//if err != nil || fetchedTask == nil {
	//	return domain.StatusResponse{
	//		StatusCode: domain.NotFound,
	//	}, nil
	//}
	//
	//err = i.challengesRepo.Update(task)
	//if err != nil {
	//	return domain.StatusResponse{}, errors.Wrap(err, "cannot update task")
	//}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *ChallengesUseCase) Remove(userId, taskId int32) (domain.StatusResponse, error) {

	//task, err := i.challengesRepo.FetchById(taskId)
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
	//err = i.challengesRepo.Remove(task.Id)
	//if err != nil {
	//	return domain.StatusResponse{}, errors.Wrap(err, "cannot remove task")
	//}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}
