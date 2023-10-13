package delivery

import (
	"context"
	"microservice/app"
	"microservice/app/core"
	"microservice/domain"
	pb "microservice/pkg/pb/api"
)

type StatusDeliveryService struct {
	pb.StatusServiceServer
	log core.Logger
}

func NewStatusDeliveryService(log core.Logger) *StatusDeliveryService {
	return &StatusDeliveryService{
		log: log,
	}
}

func (d *StatusDeliveryService) Init() error {
	app.InitGRPCService(pb.RegisterStatusServiceServer, pb.StatusServiceServer(d))
	return nil
}

func (d *StatusDeliveryService) Ping(ctx context.Context, r *pb.EmptyMessage) (*pb.StatusResponse, error) {
	return &pb.StatusResponse{
		Status: &pb.Status{
			Code:    domain.Success,
			Message: domain.Success,
		},
	}, nil
}
