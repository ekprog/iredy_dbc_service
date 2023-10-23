package usecase

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
	"microservice/layers/services"
	"microservice/tools"
	"strings"
	"time"
)

type ChallengesUseCase struct {
	log                 core.Logger
	usersRepo           domain.UsersRepository
	categoryRepo        domain.DBCCategoryRepository
	challengesRepo      domain.DBCChallengesRepository
	tracksRepo          domain.DBCTrackRepository
	periodTypeGenerator *services.PeriodTypeProcessor
	trackProcessor      *services.DBCTrackProcessor
}

func NewChallengesUseCase(log core.Logger,
	usersRepo domain.UsersRepository,
	projectsRepo domain.DBCCategoryRepository,
	periodTypeGenerator *services.PeriodTypeProcessor,
	tasksRepo domain.DBCChallengesRepository,
	tracksRepo domain.DBCTrackRepository,
	trackProcessor *services.DBCTrackProcessor) *ChallengesUseCase {
	return &ChallengesUseCase{
		log:                 log,
		usersRepo:           usersRepo,
		categoryRepo:        projectsRepo,
		challengesRepo:      tasksRepo,
		tracksRepo:          tracksRepo,
		periodTypeGenerator: periodTypeGenerator,
		trackProcessor:      trackProcessor,
	}
}

// Returns all challenges of user with some last tracks (successful or not)
func (ucase *ChallengesUseCase) All(userId int32) (domain.ChallengesListResponse, error) {
	var items []*domain.DBCChallenge
	var err error

	// Here we get only success tracks
	items, err = ucase.challengesRepo.FetchUsersAll(userId)
	if err != nil {
		return domain.ChallengesListResponse{}, errors.Wrap(err, "cannot fetch dbc-challenges by user id")
	}

	// Добавляем к Активным челленжам последние 3 трека
	for _, item := range items {

		if item.IsAutoTrack {
			continue
		}

		// ToDo: Период генерации треков (на данный момент он равен 1 суток без возможности изменения)
		period := domain.PeriodTypeEveryDay

		// Отскочить на 3 последних треков (учитывая период их генерации)
		startTime := tools.RoundDateTimeToDay(time.Now().UTC())

		list, err := ucase.periodTypeGenerator.BackwardList(item.CreatedAt, startTime, period, 3)
		if err != nil {
			return domain.ChallengesListResponse{}, errors.Wrap(err, "All")
		}

		tracks, err := ucase.tracksRepo.FetchForChallengeByDates(item.Id, list)
		if err != nil {
			return domain.ChallengesListResponse{}, errors.Wrap(err, "All")
		}
		item.LastTracks = tracks

		// Далее итерируемся на период дней (не забывает обрезать время при сравнении)
		// и проверяем, если ли в БД выборке трек по этому дню.
		// Если его нет, то создаем с Done = false (отсутствие трека в БД говорит о его не успешности)
		//err = ucase.periodTypeGenerator.StepBackwardForEach(item.CreatedAt, startTime, period, 3, func(currentTime time.Time) {
		//	//currentTimeFormat := currentTime.Format("02-01-2006")
		//	//println(currentTimeFormat)
		//
		//	// Проверяем, есть ли трек в БД (если их нет, то пользователь не отмечал их)
		//	_, ok := lo.Find(item.LastTracks, func(x *domain.DBCTrack) bool {
		//		return tools.IsEqualDateTimeByDay(x.Date, currentTime)
		//	})
		//	if !ok {
		//		item.LastTracks = append(item.LastTracks, &domain.DBCTrack{
		//			Date: currentTime,
		//			Done: false,
		//		})
		//	}
		//})
		//if err != nil {
		//	return domain.ChallengesListResponse{}, errors.Wrap(err, "PeriodTypeProcessor")
		//}
		//
		//// Сортируем по времени
		//sort.Slice(item.LastTracks, func(i, j int) bool {
		//	return item.LastTracks[i].Date.Before(item.LastTracks[j].Date)
		//})
	}

	return domain.ChallengesListResponse{
		StatusCode: domain.Success,
		Challenges: items,
	}, nil
}

func (ucase *ChallengesUseCase) Create(form *domain.CreateDBCChallengeForm) (domain.CreateChallengeResponse, error) {

	//
	err := ucase.usersRepo.InsertIfNotExists(&domain.User{Id: form.UserId})
	if err != nil {
		return domain.CreateChallengeResponse{}, errors.Wrap(err, "Create")
	}

	// Is challenge connected to category?
	var categoryId *int32
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
		categoryId = &category.Id
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
		UserId:      form.UserId,
		CategoryId:  categoryId,
		Name:        form.Name,
		Desc:        form.Desc,
		IsAutoTrack: form.IsAutoTrack,
		LastSeries:  0,
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

func (ucase *ChallengesUseCase) TrackDay(ctx context.Context, form *domain.DBCTrack) (domain.UserGamifyResponse, error) {

	// DRY
	makeErr := func(err error) (domain.UserGamifyResponse, error) {
		return domain.UserGamifyResponse{}, errors.Wrap(err, "TrackDay")
	}

	// Пользователя не проверяем, потому что наличие челленджа уже гарантирует наличие пользователя
	challenge, err := ucase.challengesRepo.FetchById(form.ChallengeId)
	if err != nil {
		return makeErr(err)
	}

	// User is owner
	if challenge == nil || challenge.UserId != form.UserId {
		return domain.UserGamifyResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	// Check if track possible to change
	// ToDo: Make 3 step back and check dates

	// 1. Проверяем есть ли такой трек
	prevDone, err := ucase.tracksRepo.CheckDoneByDate(ctx, challenge.Id, form.Date)
	if err != nil {
		return makeErr(err)
	}

	// Изменение не имеет смысла - те же значения
	if form.Done == prevDone {
		return domain.UserGamifyResponse{
			StatusCode: domain.UserLogicError,
		}, nil
	}

	// Рассчитываем изменение ScoreDaily (Score будет меняться в горутине каждые сутки)
	scoreChange := int32(0)
	if form.Done && !prevDone {
		scoreChange = 1
	} else if !form.Done && prevDone {
		scoreChange = -1
	}

	// 2. Обновляем трек в БД
	err = ucase.tracksRepo.InsertOrUpdate(ctx, form)
	if err != nil {
		return makeErr(err)
	}

	// Gamify update
	// При изменении трека в любом случае меняется ScoreDaily
	user, err := ucase.usersRepo.FetchById(form.UserId)
	if err != nil {
		return makeErr(err)
	}

	// Обновляем ScoreDaily у пользователя
	user.ScoreDaily += scoreChange
	err = ucase.usersRepo.Update(user)
	if err != nil {
		return makeErr(err)
	}

	return domain.UserGamifyResponse{
		StatusCode: domain.Success,
		ScoreDaily: user.ScoreDaily,
		LastSeries: 0,
	}, nil
}
