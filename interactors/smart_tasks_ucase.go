package interactors

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"microservice/app/core"
	"microservice/domain"
	"time"
)

type SmartTasksInteractor struct {
	log                  core.Logger
	periodMatcherService domain.PeriodMatcher
	usersRepo            domain.UsersRepository
	projectsRepo         domain.ProjectsRepository
	tasksUCase           domain.TasksInteractor
	smartTasksRepo       domain.SmartTasksRepository
}

func NewSmartTasksInteractor(
	log core.Logger,
	periodMatcherService domain.PeriodMatcher,
	usersRepo domain.UsersRepository,
	projectsRepo domain.ProjectsRepository,
	tasksUCase domain.TasksInteractor,
	smartTasksRepo domain.SmartTasksRepository) *SmartTasksInteractor {
	return &SmartTasksInteractor{
		log:                  log,
		periodMatcherService: periodMatcherService,
		usersRepo:            usersRepo,
		projectsRepo:         projectsRepo,
		tasksUCase:           tasksUCase,
		smartTasksRepo:       smartTasksRepo,
	}
}

func (i *SmartTasksInteractor) All(userId int32, trashed bool, offset, limit int32) (domain.SmartTasksListResponse, error) {

	items, err := i.smartTasksRepo.FetchAllByUserId(userId, trashed, offset, limit)
	if err != nil {
		return domain.SmartTasksListResponse{}, errors.Wrap(err, "cannot fetch smart tasks by user id")
	}

	return domain.SmartTasksListResponse{
		StatusCode: domain.Success,
		Tasks:      items,
	}, nil
}

func (i *SmartTasksInteractor) GenerateTasks() (domain.StatusResponse, error) {

	// Calculate count of tasks
	count, err := i.smartTasksRepo.CountAllActive()
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "error while counting all active smart tasks")
	}
	log.Infof("Was found %d smart tasks that required generation user tasks", count)

	// Split into pages for chunks handling
	perPage := int32(2)
	pagesCount := count / perPage
	if count%perPage > 0 {
		pagesCount++
	}

	// Handle each page
	offset := int32(0)
	for k := int32(0); k < pagesCount; k++ {
		// Load smart tasks for page
		smartTasks, err := i.smartTasksRepo.FetchAllActive(offset, perPage)
		if err != nil {
			return domain.StatusResponse{}, errors.Wrap(err, "error while fetching all active smart tasks")
		}
		offset += perPage

		// Generation for each smart task
		for _, smartTask := range smartTasks {
			i.log.Debug("SMART TASK id = %d", smartTask.Id)

			// Размер шага = 1 день
			// Начинаем с последнего дня, когда была генерация для этой умной задачи
			// Если генерации ни разу не было, то ставим в последнюю дату = CreatedAt - 1 день
			var timeIter time.Time
			if smartTask.LastGeneratedAt == nil {
				timeIter = smartTask.CreatedAt.Add(-24 * time.Hour)
			} else {
				timeIter = *smartTask.LastGeneratedAt
			}

			// Начало
			timeIter = timeIter.Add(24 * time.Hour)
			timeTo := time.Now()

			// Для всего времени генерации с шагом в сутки
			for timeTo.Truncate(24 * time.Hour).Equal(timeIter.Truncate(24 * time.Hour)) {
				for _, gItem := range smartTask.GenerationItems {
					i.log.Debug("Item period is %s", gItem.Period)
					i.log.Debug("Iter = %s", timeIter.String())
					i.log.Debug("Template = %s", gItem.Datetime.String())

					match, err := i.periodMatcherService.Match(timeIter, gItem.Datetime, gItem.Period)
					if err != nil {
						return domain.StatusResponse{}, errors.Wrap(err, "error while match smart task's generation item")
					}

					i.log.Debug("Item match = %t", match)
					if match {
						// Creating task
						createStatus, err := i.tasksUCase.Create(&domain.Task{
							UserId:    smartTask.UserId,
							ProjectId: smartTask.ProjectId,
							Name:      timeIter.String(),
							Desc:      smartTask.Desc,
							Priority:  smartTask.Priority,
						})
						if err != nil {
							panic(err)
						}
						if createStatus.StatusCode != domain.Success {
							panic("not success")
						}
					}

				}
				timeIter = timeIter.Add(24 * time.Hour)
			}

			// update smart task
			err := i.smartTasksRepo.UpdateLastGeneratedAt(smartTask.Id, timeIter.Add(-24*time.Hour))
			if err != nil {
				panic(err)
			}
		}
	}

	println(pagesCount)
	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}
