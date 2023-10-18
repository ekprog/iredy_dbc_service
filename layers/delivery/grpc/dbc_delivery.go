package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"microservice/app"
	"microservice/app/conv"
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

func (d *DBCDeliveryService) CreateChallenge(ctx context.Context, r *pb.CreateChallengeRequest) (*pb.CreateChallengesResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.dbcChallengesUCase.Create(&domain.CreateDBCChallengeForm{
		UserId:       userId,
		Name:         r.Name,
		CategoryName: r.CategoryName,
		Desc:         conv.ValueOrDefault(r.Desc),
	})
	if err != nil {
		return nil, errors.Wrap(err, "CreateChallenge")
	}

	response := &pb.CreateChallengesResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.Id = uCaseRes.Id
		response.CategoryId = uCaseRes.CategoryId
	}

	return response, nil
}

func (d *DBCDeliveryService) UpdateCategory(ctx context.Context, r *pb.UpdateCategoriesRequest) (*pb.StatusResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	category := &domain.DBCCategory{
		Id:     r.Id,
		UserId: userId,
		Name:   r.Name,
	}

	uCaseRes, err := d.dbcCategoriesUCase.Update(category)
	if err != nil {
		return nil, errors.Wrap(err, "UpdateCategory")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}
	return response, nil
}

func (d *DBCDeliveryService) GetChallenges(ctx context.Context, r *pb.GetChallengesRequest) (*pb.GetChallengesResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.dbcChallengesUCase.All(userId)
	if err != nil {
		return nil, errors.Wrap(err, "GetChallenges")
	}

	response := &pb.GetChallengesResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		Challenges: []*pb.DBCChallenge{},
	}

	if uCaseRes.StatusCode == domain.Success {
		for _, pItem := range uCaseRes.Challenges {
			p := &pb.DBCChallenge{
				Id:         pItem.Id,
				UserId:     pItem.UserId,
				Name:       pItem.Name,
				CategoryId: pItem.CategoryId,
				CreatedAt:  timestamppb.New(pItem.CreatedAt),
				DeletedAt:  conv.NullableTime(pItem.DeletedAt),
				UpdatedAt:  timestamppb.New(pItem.UpdatedAt),
				LastTracks: []*pb.DBTrack{},
			}

			for _, pTrack := range pItem.LastTracks {
				t := &pb.DBTrack{
					Date:       timestamppb.New(pTrack.Date),
					DateString: pTrack.Date.Format("02-01-2006"),
					Done:       pTrack.Done,
				}
				p.LastTracks = append(p.LastTracks, t)
			}
			response.Challenges = append(response.Challenges, p)
		}
	}

	return response, nil
}

func (d *DBCDeliveryService) GetCategories(ctx context.Context, r *pb.GetCategoriesRequest) (*pb.GetCategoriesResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}
	uCaseRes, err := d.dbcCategoriesUCase.Get(userId)
	if err != nil {
		return nil, errors.Wrap(err, "cannot fetch project info")
	}

	response := &pb.GetCategoriesResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success && uCaseRes.Categories != nil {
		for _, category := range uCaseRes.Categories {
			p := &pb.DBCCategory{
				UserId:    category.UserId,
				Id:        category.Id,
				Name:      category.Name,
				CreatedAt: timestamppb.New(category.CreatedAt),
				UpdatedAt: timestamppb.New(category.UpdatedAt),
			}

			if category.DeletedAt != nil {
				p.DeletedAt = timestamppb.New(*category.DeletedAt)
			}

			response.Categories = append(response.Categories, p)
		}
	}

	return response, nil
}

func (d *DBCDeliveryService) RemoveCategory(ctx context.Context, r *pb.RemoveCategoriesRequest) (*pb.StatusResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCase, err := d.dbcCategoriesUCase.Remove(userId, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "cannot remove category ")
	}

	response := &pb.StatusResponse{
		Status: &pb.Status{
			Code:    uCase.StatusCode,
			Message: uCase.StatusCode,
		},
	}
	return response, nil
}
