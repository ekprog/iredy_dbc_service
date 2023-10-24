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

type UsersDeliveryService struct {
	pb.UsersServiceServer
	log        core.Logger
	usersUCase domain.UsersUseCase
}

func NewUsersDeliveryService(log core.Logger,
	usersUCase domain.UsersUseCase) *UsersDeliveryService {
	return &UsersDeliveryService{
		log:        log,
		usersUCase: usersUCase,
	}
}

func (d *UsersDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterUsersServiceServer, pb.UsersServiceServer(d))
	return nil
}

func (d *UsersDeliveryService) Info(ctx context.Context, r *pb.IdRequest) (*pb.GetUserResponse, error) {
	uCaseRes, err := d.usersUCase.Info(ctx, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Info")
	}

	response := &pb.GetUserResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		User: nil,
	}

	if uCaseRes.StatusCode == domain.Success {
		response.User = &pb.User{
			Id:         uCaseRes.User.Id,
			Score:      uCaseRes.User.Score,
			ScoreDaily: uCaseRes.User.ScoreDaily,
			CreatedAt:  timestamppb.New(uCaseRes.User.CreatedAt),
			UpdatedAt:  timestamppb.New(uCaseRes.User.UpdatedAt),
			DeletedAt:  conv.NullableTime(uCaseRes.User.DeletedAt),
		}
	}

	return response, nil
}

func (d *UsersDeliveryService) MyInfo(ctx context.Context, r *pb.EmptyMessage) (*pb.GetUserResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.usersUCase.Info(ctx, userId)
	if err != nil {
		return nil, errors.Wrap(err, "MyInfo")
	}

	response := &pb.GetUserResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		User: nil,
	}

	if uCaseRes.StatusCode == domain.Success {
		response.User = &pb.User{
			Id:         uCaseRes.User.Id,
			Score:      uCaseRes.User.Score,
			ScoreDaily: uCaseRes.User.ScoreDaily,
			CreatedAt:  timestamppb.New(uCaseRes.User.CreatedAt),
			UpdatedAt:  timestamppb.New(uCaseRes.User.UpdatedAt),
			DeletedAt:  conv.NullableTime(uCaseRes.User.DeletedAt),
		}
	}

	return response, nil
}
