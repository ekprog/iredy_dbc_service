package interactors

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/domain"
)

type ProjectsInteractor struct {
	log          core.Logger
	projectsRepo domain.ProjectsRepository
	usersRepo    domain.UsersRepository
	tasksRepo    domain.TasksRepository
}

func NewProjectsInteractor(log core.Logger,
	usersRepo domain.UsersRepository,
	projectsRepo domain.ProjectsRepository,
	tasksRepo domain.TasksRepository) *ProjectsInteractor {
	return &ProjectsInteractor{
		log:          log,
		usersRepo:    usersRepo,
		projectsRepo: projectsRepo,
		tasksRepo:    tasksRepo,
	}
}

func (i *ProjectsInteractor) Active(userId int32, trashed bool) (domain.ProjectListResponse, error) {

	var projects []*domain.Project
	var err error

	if trashed {
		projects, err = i.projectsRepo.FetchByUserIdTrashed(userId)
		if err != nil {
			return domain.ProjectListResponse{}, errors.Wrap(err, "cannot fetch active projects by user id")
		}
	} else {
		projects, err = i.projectsRepo.FetchByUserId(userId)
		if err != nil {
			return domain.ProjectListResponse{}, errors.Wrap(err, "cannot fetch trashed projects by user id")
		}
	}

	// Подгружаем задачи, если необходимо
	// ToDo: Агрегировать в один запрос
	//if withTasks {
	//	for _, p := range projects {
	//		// Active Tasks
	//		tasks, err := i.tasksRepo.FetchByProjectId(p.Id, false, 0, 50)
	//		if err != nil {
	//			return domain.ProjectListResponse{}, errors.Wrapf(err, "cannot fetch projects (%d) active tasks", p.Id)
	//		}
	//		p.Tasks = tasks
	//
	//		// UnDone Tasks
	//		doneTasks, err := i.tasksRepo.FetchByProjectId(p.Id, true, 0, 50)
	//		if err != nil {
	//			return domain.ProjectListResponse{}, errors.Wrapf(err, "cannot fetch projects (%d) done tasks", p.Id)
	//		}
	//		p.DoneTasks = doneTasks
	//	}
	//}

	return domain.ProjectListResponse{
		StatusCode: domain.Success,
		Projects:   projects,
	}, nil
}

func (i *ProjectsInteractor) Info(userId int32, projectId int32) (domain.ProjectInfoResponse, error) {

	project, err := i.projectsRepo.FetchById(projectId)
	if err != nil {
		return domain.ProjectInfoResponse{}, errors.Wrapf(err, "cannot fetch project by id %d", projectId)
	}

	if project == nil {
		return domain.ProjectInfoResponse{
			StatusCode: domain.ProjectNotFound,
		}, err
	}

	if project.UserId != userId {
		return domain.ProjectInfoResponse{
			StatusCode: domain.AccessDenied,
		}, nil
	}

	return domain.ProjectInfoResponse{
		StatusCode: domain.Success,
		Project:    project,
	}, nil
}

func (i *ProjectsInteractor) Trashed(userId int32) (domain.ProjectListResponse, error) {
	projects, err := i.projectsRepo.FetchByUserIdTrashed(userId)
	if err != nil {
		return domain.ProjectListResponse{}, errors.Wrap(err, "cannot fetch trashed projects by user id")
	}

	return domain.ProjectListResponse{
		StatusCode: domain.Success,
		Projects:   projects,
	}, nil
}

func (i *ProjectsInteractor) Create(project domain.Project) (domain.IdResponse, error) {

	if project.Name == "" || project.UserId < 0 {
		return domain.IdResponse{
			StatusCode: domain.ValidationError,
		}, nil
	}

	// If user does not exist - create
	err := i.usersRepo.InsertIfNotExists(&domain.User{
		Id: project.UserId,
	})
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot insert user before creating project")
	}

	// CreateIfNotExists project
	err = i.projectsRepo.Insert(&project)
	if err != nil {
		return domain.IdResponse{}, errors.Wrap(err, "cannot insert project")
	}

	return domain.IdResponse{
		StatusCode: domain.Success,
		Id:         project.Id,
	}, nil
}

func (i *ProjectsInteractor) Remove(userId, projectId int32) (domain.StatusResponse, error) {

	// Check if user's project exists
	project, err := i.projectsRepo.FetchById(projectId)
	if err != nil {
		return domain.StatusResponse{},
			errors.Wrapf(err, "cannot fetch project before removing. ProjectId=%d", projectId)
	}

	if project.UserId != userId {
		return domain.StatusResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	// Remove
	err = i.projectsRepo.Remove(projectId)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot remove project %d", projectId)
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *ProjectsInteractor) Update(item *domain.Project) (domain.StatusResponse, error) {

	// Check if user's project exists
	project, err := i.projectsRepo.FetchById(item.Id)
	if err != nil {
		return domain.StatusResponse{},
			errors.Wrapf(err, "cannot fetch project before updating. ProjectId=%d", item.Id)
	}

	if project.UserId != item.UserId {
		return domain.StatusResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	// Update
	err = i.projectsRepo.Update(project)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot update project %d", project.Id)
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *ProjectsInteractor) UpdateOrder(userId int32, newOrder []int32) (domain.StatusResponse, error) {

	// Todo: check all projects belong to user (now it is in sql)
	//...

	// Update for user
	err := i.projectsRepo.UpdateOrderForUser(userId, newOrder)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot update order for user projects")
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *ProjectsInteractor) DragProjects(userId int32, drags []*domain.ProjectDrag) (domain.StatusResponse, error) {

	// Todo: check all projects belong to user (now it is in sql)
	//...

	// Update for user
	err := i.projectsRepo.DragItemsForUser(userId, drags)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot drag projects for user %d", userId)
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}
