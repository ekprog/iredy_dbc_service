package grpc

import (
	"context"
	"github.com/pkg/errors"
	"microservice/app"
	"microservice/app/core"
	"microservice/layers/domain"
	pb "microservice/pkg/pb/api"
)

type DBCDeliveryService struct {
	pb.DBCServiceServer
	log                core.Logger
	usersUCase         domain.UsersUseCase
	dbcCategoriesUCase domain.DBCCategoryUseCase
	dbcChallengesUCase domain.ChallengesUseCase
}

func NewDBCDeliveryService(log core.Logger,
	usersUCase domain.UsersUseCase,
	dbcCategoriesUCase domain.DBCCategoryUseCase,
	dbcChallengesUCase domain.ChallengesUseCase) *DBCDeliveryService {
	return &DBCDeliveryService{
		log:                log,
		usersUCase:         usersUCase,
		dbcCategoriesUCase: dbcCategoriesUCase,
		dbcChallengesUCase: dbcChallengesUCase,
	}
}

func (d *DBCDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterDBCServiceServer, pb.DBCServiceServer(d))
	return nil
}

func (d *DBCDeliveryService) CreateChallenge(ctx context.Context, r *pb.CreateChallengeRequest) (*pb.IdResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	// Есть ли такой пользоваетль?
	// Если нет - создать пользователя + создать запись

	uCaseRes, err := d.dbcChallengesUCase.Create(&domain.CreateDBCChallengeForm{
		UserId: userId,
		Name:   r.Name,
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

//func (d *DBCDeliveryService) GetProjects(ctx context.Context, r *pb.GetProjectsRequest) (*pb.GetProjectsResponse, error) {
//	userId, err := app.ExtractRequestUserId(ctx)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot extract user_id from context")
//	}
//
//	uCaseRes, err := d.projectsUCase.Active(
//		userId,
//		conv.ValueOrDefault(r.Trashed, false),
//	)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot fetch projects")
//	}
//
//	response := &pb.GetProjectsResponse{
//		Status: &pb.Status{
//			Code:    uCaseRes.StatusCode,
//			Message: uCaseRes.StatusCode,
//		},
//		Projects: []*pb.Project{},
//	}
//
//	if uCaseRes.StatusCode == domain.Success && uCaseRes.Projects != nil {
//		for _, pItem := range uCaseRes.Projects {
//			p := &pb.Project{
//				Id:        pItem.Id,
//				UserId:    pItem.UserId,
//				Name:      pItem.Name,
//				Desc:      pItem.Desc,
//				Color:     pItem.Color,
//				Order:     pItem.Order,
//				ParentId:  pItem.ParentId,
//				CreatedAt: timestamppb.New(pItem.CreatedAt),
//				UpdatedAt: timestamppb.New(pItem.UpdatedAt),
//			}
//			if pItem.DeletedAt != nil {
//				p.DeletedAt = timestamppb.New(*pItem.DeletedAt)
//			}
//
//			//if pItem.Challenges != nil {
//			//	for _, tItem := range pItem.Challenges {
//			//		t := &pb.DBCChallenge{
//			//			Id:        tItem.Id,
//			//			UserId:    tItem.UserId,
//			//			CategoryId: tItem.CategoryId,
//			//			Name:      tItem.Name,
//			//			Desc:      tItem.Desc,
//			//			Priority:  int32(tItem.Priority),
//			//			Done:      tItem.Done,
//			//			CreatedAt: timestamppb.New(tItem.CreatedAt),
//			//			UpdatedAt: timestamppb.New(tItem.UpdatedAt),
//			//		}
//			//		if tItem.DeletedAt != nil {
//			//			t.DeletedAt = timestamppb.New(*tItem.DeletedAt)
//			//		}
//			//		p.Challenges = append(p.Challenges, t)
//			//	}
//			//}
//			//
//			//if pItem.DoneTasks != nil {
//			//	for _, tItem := range pItem.DoneTasks {
//			//		t := &pb.DBCChallenge{
//			//			Id:        tItem.Id,
//			//			UserId:    tItem.UserId,
//			//			CategoryId: tItem.CategoryId,
//			//			Name:      tItem.Name,
//			//			Desc:      tItem.Desc,
//			//			Priority:  int32(tItem.Priority),
//			//			Done:      tItem.Done,
//			//			CreatedAt: timestamppb.New(tItem.CreatedAt),
//			//			UpdatedAt: timestamppb.New(tItem.UpdatedAt),
//			//		}
//			//		if tItem.DeletedAt != nil {
//			//			t.DeletedAt = timestamppb.New(*tItem.DeletedAt)
//			//		}
//			//		p.HistoryTasks = append(p.HistoryTasks, t)
//			//	}
//			//}
//
//			response.Projects = append(response.Projects, p)
//		}
//	}
//
//	return response, nil
//}
//
//func (d *DBCDeliveryService) GetProjectInfo(ctx context.Context, r *pb.IdRequest) (*pb.GetProjectInfoResponse, error) {
//	userId, err := app.ExtractRequestUserId(ctx)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot extract user_id from context")
//	}
//
//	uCaseRes, err := d.projectsUCase.Info(userId, r.Id)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot fetch project info")
//	}
//
//	response := &pb.GetProjectInfoResponse{
//		Status: &pb.Status{
//			Code:    uCaseRes.StatusCode,
//			Message: uCaseRes.StatusCode,
//		},
//	}
//
//	if uCaseRes.StatusCode == domain.Success && uCaseRes.Project != nil {
//		response.Project = &pb.Project{
//			Id:        uCaseRes.Project.Id,
//			UserId:    uCaseRes.Project.UserId,
//			Name:      uCaseRes.Project.Name,
//			Desc:      uCaseRes.Project.Desc,
//			Color:     uCaseRes.Project.Color,
//			CreatedAt: timestamppb.New(uCaseRes.Project.CreatedAt),
//			UpdatedAt: timestamppb.New(uCaseRes.Project.UpdatedAt),
//		}
//		if uCaseRes.Project.DeletedAt != nil {
//			response.Project.DeletedAt = timestamppb.New(*uCaseRes.Project.DeletedAt)
//		}
//	}
//
//	return response, nil
//}
//
//func (d *DBCDeliveryService) RemoveProject(ctx context.Context, r *pb.IdRequest) (*pb.StatusResponse, error) {
//	userId, err := app.ExtractRequestUserId(ctx)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot extract user_id from context")
//	}
//
//	uCaseRes, err := d.projectsUCase.Remove(userId, r.Id)
//	if err != nil {
//		return nil, errors.Wrap(err, "cannot remove project")
//	}
//
//	response := &pb.StatusResponse{
//		Status: &pb.Status{
//			Code:    uCaseRes.StatusCode,
//			Message: uCaseRes.StatusCode,
//		},
//	}
//
//	return response, nil
//}
