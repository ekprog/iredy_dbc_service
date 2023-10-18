package usecase

import (
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/layers/services"
	"sort"
	"strings"
	"time"
)

type ChallengesUseCase struct {
	log                 core.Logger
	usersRepo           domain.UsersRepository
	categoryRepo        domain.DBCCategoryRepository
	challengesRepo      domain.DBCChallengesRepository
	usersUseCase        domain.UsersUseCase
	periodTypeGenerator *services.PeriodTypeGenerator
}

func NewChallengesUseCase(log core.Logger,
	usersRepo domain.UsersRepository,
	projectsRepo domain.DBCCategoryRepository,
	periodTypeGenerator *services.PeriodTypeGenerator,
	tasksRepo domain.DBCChallengesRepository,
	usersUseCase domain.UsersUseCase) *ChallengesUseCase {
	return &ChallengesUseCase{
		log:                 log,
		usersRepo:           usersRepo,
		categoryRepo:        projectsRepo,
		challengesRepo:      tasksRepo,
		usersUseCase:        usersUseCase,
		periodTypeGenerator: periodTypeGenerator,
	}
}

// Returns all challenges of user with some last tracks (successful or not)
func (ucase *ChallengesUseCase) All(userId int32) (domain.ChallengesListResponse, error) {
	var items []*domain.DBCChallenge
	var err error

	// Here we get only success tracks
	items, err = ucase.challengesRepo.FetchAll(userId)
	if err != nil {
		return domain.ChallengesListResponse{}, errors.Wrap(err, "cannot fetch dbc-challenges by user id")
	}

	for _, item := range items {
		// ToDo: Период генерации треков (на данный момент он равен 1 суток без возможности изменения)
		period := domain.PeriodTypeEveryDay

		// Отскочить на 5 последних треков (учитывая период их генерации)
		startTime, err := ucase.periodTypeGenerator.Step(item.CreatedAt, time.Now(), period, -5)
		if err != nil {
			return domain.ChallengesListResponse{}, errors.Wrap(err, "PeriodTypeGenerator")
		}

		// Далее итерируемся на период дней (не забывает обрезать время при сравнении)
		// и проверяем, если ли в БД выборке трек по этому дню.
		// Если его нет, то создаем с Done = false (отсутствие трека в БД говорит о его не успешности)
		err = ucase.periodTypeGenerator.StepForwardForEach(item.CreatedAt, startTime, period, 5, func(currentTime time.Time) {
			currentTimeFormat := currentTime.Format("02-01-2016")
			println(currentTimeFormat)

			// Проверяем, есть ли трек в БД (если их нет, то пользователь не отмечал их)
			_, ok := lo.Find(item.LastTracks, func(x *domain.DBCTrack) bool {
				return x.Date.Format("02-01-2016") == currentTimeFormat
			})
			if !ok {
				item.LastTracks = append(item.LastTracks, &domain.DBCTrack{
					Date: currentTime,
					Done: false,
				})
			}
		})
		if err != nil {
			return domain.ChallengesListResponse{}, errors.Wrap(err, "PeriodTypeGenerator")
		}

		// Сортируем по времени
		sort.Slice(item.LastTracks, func(i, j int) bool {
			return item.LastTracks[i].Date.Before(item.LastTracks[j].Date)
		})
	}

	return domain.ChallengesListResponse{
		StatusCode: domain.Success,
		Challenges: items,
	}, nil
}

func (ucase *ChallengesUseCase) Create(form *domain.CreateDBCChallengeForm) (domain.CreateChallengeResponse, error) {

	//
	_, err := ucase.usersUseCase.CreateIfNotExists(domain.User{Id: form.UserId})
	if err != nil {
		return domain.CreateChallengeResponse{}, errors.Wrap(err, "cannot insert user before creating task")
	}

	// Is challenge connected to category?
	categoryId := new(int32)
	if form.CategoryName != nil {
		// Finding category
		category, err := ucase.categoryRepo.FetchByName(form.UserId, *form.CategoryName)
		if err != nil {
			return domain.CreateChallengeResponse{}, errors.Wrap(err, "cannot fetch category before creating task")
		}

		// Creating if not exists
		if category == nil {
			category = &domain.DBCCategory{
				UserId: form.UserId,
				Name:   *form.CategoryName,
			}
			err = ucase.categoryRepo.Insert(category)
			if err != nil {
				return domain.CreateChallengeResponse{}, errors.Wrap(err, "cannot insert new category before creating task")
			}
		}

		// Set category Id for next step
		*categoryId = category.Id
	}

	// Check if challenge with same name already exists
	form.Name = strings.TrimSpace(form.Name)
	challengeFound, err := ucase.challengesRepo.FetchByName(form.UserId, form.Name)
	if err != nil {
		return domain.CreateChallengeResponse{}, errors.Wrap(err, "cannot check if challenge exists by name before creating task")
	}
	if challengeFound != nil {
		return domain.CreateChallengeResponse{
			StatusCode: domain.AlreadyExists,
		}, nil
	}

	// Validation of challenge form
	if form.Name == "" {
		return domain.CreateChallengeResponse{
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
		return domain.CreateChallengeResponse{}, errors.Wrap(err, "cannot insert new challenge before creating task")
	}

	//
	return domain.CreateChallengeResponse{
		StatusCode: domain.Success,
		Id:         challenge.Id,
		CategoryId: categoryId,
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
