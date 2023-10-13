package interactors

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
	"microservice/kafka"
	"microservice/pkg/pb/events"
)

type TasksInteractor struct {
	log          core.Logger
	usersRepo    domain.UsersRepository
	projectsRepo domain.ProjectsRepository
	tasksRepo    domain.TasksRepository
}

func NewTasksInteractor(log core.Logger,
	usersRepo domain.UsersRepository,
	projectsRepo domain.ProjectsRepository,
	tasksRepo domain.TasksRepository) *TasksInteractor {
	return &TasksInteractor{log: log, usersRepo: usersRepo, projectsRepo: projectsRepo, tasksRepo: tasksRepo}
}

func (i *TasksInteractor) All(userId int32, projectId *int32, done bool, offset, limit int32) (domain.TasksListResponse, error) {
	var items []*domain.Task
	var err error

	// If task is a child of project
	if projectId != nil {
		project, err := i.projectsRepo.FetchById(*projectId)
		if err != nil {
			return domain.TasksListResponse{}, errors.Wrapf(err, "cannot fetch project %d", *projectId)
		}
		if project == nil {
			return domain.TasksListResponse{
				StatusCode: domain.ProjectNotFound,
			}, nil
		}
		// Is user owner of project?
		if project.UserId != userId {
			return domain.TasksListResponse{
				StatusCode: domain.AccessDenied,
			}, nil
		}
	}

	if projectId == nil {
		items, err = i.tasksRepo.FetchAll(userId, done, offset, limit)
		if err != nil {
			return domain.TasksListResponse{}, errors.Wrap(err, "cannot fetch tasks by user id")
		}
	} else {
		items, err = i.tasksRepo.FetchByProjectId(*projectId, done, offset, limit)
		if err != nil {
			return domain.TasksListResponse{}, errors.Wrapf(err, "cannot fetch project tasks by project id %d", *projectId)
		}
	}

	return domain.TasksListResponse{
		StatusCode: domain.Success,
		Tasks:      items,
	}, nil
}

func (i *TasksInteractor) Create(task *domain.Task) (domain.IdResponse, error) {

	if task.Name == "" || task.Priority < 0 {
		return domain.IdResponse{
			StatusCode: domain.ValidationError,
		}, nil
	}

	// If user does not exist - create
	err := i.usersRepo.InsertIfNotExists(&domain.User{
		Id: task.UserId,
	})
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot insert user before creating task")
	}

	// If task is a child of project
	if task.ProjectId != nil {
		project, err := i.projectsRepo.FetchById(*task.ProjectId)
		if err != nil {
			return domain.IdResponse{}, errors.Wrapf(err, "cannot fetch project %d", *task.ProjectId)
		}
		if project == nil {
			return domain.IdResponse{
				StatusCode: domain.ProjectNotFound,
			}, nil
		}
		// Is user owner of project?
		if project.UserId != task.UserId {
			return domain.IdResponse{
				StatusCode: domain.AccessDenied,
			}, nil
		}
	}

	err = i.tasksRepo.Insert(task)
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot insert task")
	}

	return domain.IdResponse{
		StatusCode: domain.Success,
		Id:         task.Id,
	}, nil
}

func (i *TasksInteractor) Update(task *domain.Task) (domain.StatusResponse, error) {

	fetchedTask, err := i.tasksRepo.FetchById(task.Id)
	if err != nil || fetchedTask == nil {
		return domain.StatusResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	err = i.tasksRepo.Update(task)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "cannot update task")
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *TasksInteractor) SetDone(userId, taskId int32, flag bool) (domain.StatusResponse, error) {

	task, err := i.tasksRepo.FetchById(taskId)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot fetch task by id %d", taskId)
	}

	if task == nil {
		return domain.StatusResponse{
			StatusCode: domain.TaskNotFound,
		}, nil
	}

	if task.UserId != userId {
		return domain.StatusResponse{
			StatusCode: domain.AccessDenied,
		}, nil
	}

	task.Done = flag
	err = i.tasksRepo.Update(task)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "cannot setDone task")
	}

	// Produce Event
	err = kafka.TaskDoneTopic.Produce(&events.TaskDoneEvent{
		UserId:   task.UserId,
		TaskId:   task.Id,
		Priority: int32(task.Priority),
	})
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "cannot send kafka event (TaskDoneEvent")
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *TasksInteractor) Remove(userId, taskId int32) (domain.StatusResponse, error) {

	task, err := i.tasksRepo.FetchById(taskId)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot fetch task by id %d", taskId)
	}

	if task.UserId != userId {
		return domain.StatusResponse{
			StatusCode: domain.AccessDenied,
		}, nil
	}

	err = i.tasksRepo.Remove(task.Id)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "cannot remove task")
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}
