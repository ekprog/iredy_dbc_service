package delivery

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservice/app"
	"microservice/app/conv"
	"microservice/app/core"
	"microservice/domain"
	pb "microservice/pkg/pb/api"
)

type ToDoDeliveryService struct {
	pb.UnsafeToDoServiceServer
	log             core.Logger
	usersUCase      domain.UsersInteractor
	projectsUCase   domain.ProjectsInteractor
	tasksUCase      domain.TasksInteractor
	smartTasksUCase domain.SmartTasksInteractor
}

func NewToDoDeliveryService(log core.Logger,
	usersUCase domain.UsersInteractor,
	projectsUCase domain.ProjectsInteractor,
	tasksUCase domain.TasksInteractor,
	smartTasksUCase domain.SmartTasksInteractor) *ToDoDeliveryService {
	return &ToDoDeliveryService{
		log:             log,
		usersUCase:      usersUCase,
		projectsUCase:   projectsUCase,
		tasksUCase:      tasksUCase,
		smartTasksUCase: smartTasksUCase,
	}
}

func (d *ToDoDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterToDoServiceServer, pb.ToDoServiceServer(d))
	return nil
}

func (d *ToDoDeliveryService) CreateProject(ctx context.Context, r *pb.CreateProjectRequest) (*pb.IdResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.projectsUCase.Create(domain.Project{
		UserId:   userId,
		Name:     r.Name,
		Desc:     r.Desc,
		Color:    r.Color,
		Order:    r.Order,
		ParentId: r.ParentId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create project")
	}

	response := &pb.IdResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.Id = uCaseRes.Id
	}

	return response, nil
}

func (d *ToDoDeliveryService) GetProjects(ctx context.Context, r *pb.GetProjectsRequest) (*pb.GetProjectsResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.projectsUCase.Active(
		userId,
		conv.ValueOrDefault(r.Trashed, false),
	)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch projects")
	}

	response := &pb.GetProjectsResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		Projects: []*pb.Project{},
	}

	if uCaseRes.StatusCode == domain.Success && uCaseRes.Projects != nil {
		for _, pItem := range uCaseRes.Projects {
			p := &pb.Project{
				Id:        pItem.Id,
				UserId:    pItem.UserId,
				Name:      pItem.Name,
				Desc:      pItem.Desc,
				Color:     pItem.Color,
				Order:     pItem.Order,
				ParentId:  pItem.ParentId,
				CreatedAt: timestamppb.New(pItem.CreatedAt),
				UpdatedAt: timestamppb.New(pItem.UpdatedAt),
			}
			if pItem.DeletedAt != nil {
				p.DeletedAt = timestamppb.New(*pItem.DeletedAt)
			}

			//if pItem.Tasks != nil {
			//	for _, tItem := range pItem.Tasks {
			//		t := &pb.Task{
			//			Id:        tItem.Id,
			//			UserId:    tItem.UserId,
			//			ProjectId: tItem.ProjectId,
			//			Name:      tItem.Name,
			//			Desc:      tItem.Desc,
			//			Priority:  int32(tItem.Priority),
			//			Done:      tItem.Done,
			//			CreatedAt: timestamppb.New(tItem.CreatedAt),
			//			UpdatedAt: timestamppb.New(tItem.UpdatedAt),
			//		}
			//		if tItem.DeletedAt != nil {
			//			t.DeletedAt = timestamppb.New(*tItem.DeletedAt)
			//		}
			//		p.Tasks = append(p.Tasks, t)
			//	}
			//}
			//
			//if pItem.DoneTasks != nil {
			//	for _, tItem := range pItem.DoneTasks {
			//		t := &pb.Task{
			//			Id:        tItem.Id,
			//			UserId:    tItem.UserId,
			//			ProjectId: tItem.ProjectId,
			//			Name:      tItem.Name,
			//			Desc:      tItem.Desc,
			//			Priority:  int32(tItem.Priority),
			//			Done:      tItem.Done,
			//			CreatedAt: timestamppb.New(tItem.CreatedAt),
			//			UpdatedAt: timestamppb.New(tItem.UpdatedAt),
			//		}
			//		if tItem.DeletedAt != nil {
			//			t.DeletedAt = timestamppb.New(*tItem.DeletedAt)
			//		}
			//		p.HistoryTasks = append(p.HistoryTasks, t)
			//	}
			//}

			response.Projects = append(response.Projects, p)
		}
	}

	return response, nil
}

func (d *ToDoDeliveryService) GetProjectInfo(ctx context.Context, r *pb.IdRequest) (*pb.GetProjectInfoResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.projectsUCase.Info(userId, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch project info")
	}

	response := &pb.GetProjectInfoResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success && uCaseRes.Project != nil {
		response.Project = &pb.Project{
			Id:        uCaseRes.Project.Id,
			UserId:    uCaseRes.Project.UserId,
			Name:      uCaseRes.Project.Name,
			Desc:      uCaseRes.Project.Desc,
			Color:     uCaseRes.Project.Color,
			CreatedAt: timestamppb.New(uCaseRes.Project.CreatedAt),
			UpdatedAt: timestamppb.New(uCaseRes.Project.UpdatedAt),
		}
		if uCaseRes.Project.DeletedAt != nil {
			response.Project.DeletedAt = timestamppb.New(*uCaseRes.Project.DeletedAt)
		}
	}

	return response, nil
}

func (d *ToDoDeliveryService) RemoveProject(ctx context.Context, r *pb.IdRequest) (*pb.StatusResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.projectsUCase.Remove(userId, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "cannot remove project")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) UpdateProject(ctx context.Context, r *pb.UpdateProjectRequest) (*pb.StatusResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.projectsUCase.Update(&domain.Project{
		Id:     r.ProjectId,
		UserId: userId,
		Name:   r.Name,
		Desc:   r.Desc,
		Color:  r.Color,
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot remove project")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) UpdateProjectsOrder(ctx context.Context, r *pb.UpdateProjectOrderRequest) (*pb.StatusResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	if len(r.ProjectIds) == 0 || r.ProjectIds == nil {
		return &pb.StatusResponse{
			Status: &pb.Status{
				Code:    domain.ValidationError,
				Message: "Array of ids should be not empty",
			},
		}, nil
	}

	uCaseRes, err := d.projectsUCase.UpdateOrder(userId, r.ProjectIds)
	if err != nil {
		return nil, errors.Wrap(err, "cannot update projects order")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) DragProjects(ctx context.Context, r *pb.DragProjectRequest) (*pb.StatusResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	if r.Items == nil || len(r.Items) == 0 {
		return &pb.StatusResponse{
			Status: &pb.Status{
				Code:    domain.ValidationError,
				Message: "Array of items should be not empty",
			},
		}, nil
	}

	// fill drags
	drags := make([]*domain.ProjectDrag, len(r.Items))
	for i, item := range r.Items {
		drags[i] = &domain.ProjectDrag{
			ProjectId: item.ProjectId,
			ParentId:  item.ParentId,
			Order:     item.Order,
		}
	}

	//
	uCaseRes, err := d.projectsUCase.DragProjects(userId, drags)
	if err != nil {
		return nil, errors.Wrap(err, "cannot update projects order")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) GetTasks(ctx context.Context, r *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.tasksUCase.All(userId,
		r.ProjectId,
		conv.ValueOrDefault(r.Done, false),
		conv.ValueOrDefault(r.Offset, 0),
		conv.ValueOrDefault(r.Limit, 50))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get done tasks")
	}

	response := &pb.GetTasksResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success || uCaseRes.Tasks != nil {
		for _, item := range uCaseRes.Tasks {
			p := &pb.Task{
				Id:        item.Id,
				UserId:    item.UserId,
				ProjectId: item.ProjectId,
				Name:      item.Name,
				Desc:      item.Desc,
				Priority:  int32(item.Priority),
				Done:      item.Done,
				CreatedAt: timestamppb.New(item.CreatedAt),
				UpdatedAt: timestamppb.New(item.UpdatedAt),
			}
			if item.DeletedAt != nil {
				p.DeletedAt = timestamppb.New(*item.DeletedAt)
			}
			response.Tasks = append(response.Tasks, p)
		}
	}

	return response, nil
}

func (d *ToDoDeliveryService) CreateTask(ctx context.Context, r *pb.CreateTaskRequest) (*pb.IdResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.tasksUCase.Create(&domain.Task{
		UserId:    userId,
		ProjectId: r.ProjectId,
		Name:      r.Name,
		Desc:      r.Desc,
		Priority:  domain.Priority(r.Priority),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create task")
	}

	response := &pb.IdResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.Id = uCaseRes.Id
	}

	return response, nil
}

func (d *ToDoDeliveryService) UpdateTask(ctx context.Context, r *pb.UpdateTaskRequest) (*pb.StatusResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.tasksUCase.Update(&domain.Task{
		Id:       r.TaskId,
		UserId:   userId,
		Name:     r.Name,
		Desc:     r.Desc,
		Priority: domain.Priority(r.Priority),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create task")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) SetTaskDone(ctx context.Context, r *pb.SetTaskDoneRequest) (*pb.StatusResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.tasksUCase.SetDone(userId, r.TaskId, r.Done)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot set done task %d", r.TaskId)
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) RemoveTask(ctx context.Context, r *pb.IdRequest) (*pb.StatusResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.tasksUCase.Remove(userId, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "cannot remove task")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	return response, nil
}

func (d *ToDoDeliveryService) CreateSmartTask(ctx context.Context, r *pb.CreateSmartTaskRequest) (*pb.IdResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.tasksUCase.Create(&domain.Task{
		UserId:    userId,
		ProjectId: r.ProjectId,
		Name:      r.Name,
		Desc:      r.Desc,
		Priority:  domain.Priority(r.Priority),
	})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create task")
	}

	response := &pb.IdResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.Id = uCaseRes.Id
	}

	return response, nil
}

func (d *ToDoDeliveryService) GetSmartTasks(ctx context.Context, r *pb.GetSmartTasksRequest) (*pb.GetSmartTasksResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.smartTasksUCase.All(userId,
		conv.ValueOrDefault(r.Trashed, false),
		conv.ValueOrDefault(r.Offset, 0),
		conv.ValueOrDefault(r.Limit, 1000))
	if err != nil {
		return nil, errors.Wrap(err, "cannot get smart tasks")
	}

	response := &pb.GetSmartTasksResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success && uCaseRes.Tasks != nil {
		for _, item := range uCaseRes.Tasks {

			// Task
			p := &pb.SmartTask{
				Id:        item.Id,
				UserId:    item.UserId,
				ProjectId: item.ProjectId,
				Name:      item.Name,
				Desc:      item.Desc,
				Priority:  int32(item.Priority),
				CreatedAt: timestamppb.New(item.CreatedAt),
				UpdatedAt: timestamppb.New(item.UpdatedAt),
				DeletedAt: conv.NullableTime(item.DeletedAt),
			}

			// Generation items
			for _, itemC := range item.GenerationItems {
				c := &pb.GenerationItem{
					Period:    0,
					Datetime:  timestamppb.New(itemC.Datetime),
					CreatedAt: timestamppb.New(item.CreatedAt),
					UpdatedAt: timestamppb.New(item.UpdatedAt),
					DeletedAt: conv.NullableTime(item.DeletedAt),
				}
				p.GenerationItems = append(p.GenerationItems, c)
			}

			// Append
			response.Tasks = append(response.Tasks, p)
		}
	}

	return response, nil
}
