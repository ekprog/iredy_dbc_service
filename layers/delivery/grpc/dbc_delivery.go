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
	"microservice/tools"
)

type DBCDeliveryService struct {
	pb.DBCServiceServer
	log                core.Logger
	usersUCase         domain.UsersUseCase
	dbcCategoriesUCase domain.DBCCategoryUseCase
	dbcChallengesUCase domain.DBCChallengesUseCase
}

func NewDBCDeliveryService(log core.Logger,
	usersUCase domain.UsersUseCase,
	dbcCategoriesUCase domain.DBCCategoryUseCase,
	dbcChallengesUCase domain.DBCChallengesUseCase) *DBCDeliveryService {
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

	uCaseRes, err := d.dbcChallengesUCase.UserCreate(&domain.CreateDBCChallengeForm{
		UserId:       userId,
		Name:         r.Name,
		CategoryName: r.CategoryName,
		Desc:         r.Desc,
		IsAutoTrack:  r.IsAutoTrack,
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

func (d *DBCDeliveryService) GetChallenges(ctx context.Context, r *pb.GetChallengesRequest) (*pb.GetUserChallengesResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "cannot extract user_id from context")
	}

	uCaseRes, err := d.dbcChallengesUCase.UserAll(userId)
	if err != nil {
		return nil, errors.Wrap(err, "GetChallenges")
	}

	response := &pb.GetUserChallengesResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
		Challenges: []*pb.DBCUserChallenge{},
	}

	if uCaseRes.StatusCode == domain.Success {
		for _, pItem := range uCaseRes.UserChallenges {
			p := &pb.DBCUserChallenge{
				Id:          pItem.Id,
				UserId:      pItem.UserId,
				IsAutoTrack: pItem.ChallengeInfo.IsAutoTrack,
				Name:        pItem.ChallengeInfo.Name,
				Image:       pItem.ChallengeInfo.Image,
				Desc:        pItem.ChallengeInfo.Desc,
				LastSeries:  pItem.LastSeries,
				CreatedAt:   timestamppb.New(pItem.CreatedAt),
				DeletedAt:   conv.NullableTime(pItem.DeletedAt),
				UpdatedAt:   timestamppb.New(pItem.UpdatedAt),
				LastTracks:  []*pb.DBTrack{},
			}
			if pItem.ChallengeInfo.Category != nil {
				p.CategoryId = &pItem.ChallengeInfo.Category.Id
				p.CategoryName = &pItem.ChallengeInfo.Category.Name
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

func (d *DBCDeliveryService) SearchChallenges(ctx context.Context, r *pb.SearchChallengesRequest) (*pb.GetChallengesResponse, error) {
	uCaseRes, err := d.dbcChallengesUCase.PublicSearch(r.Search, r.CategoryId, r.Limit, r.Offset)
	if err != nil {
		return nil, errors.Wrap(err, "PublicSearch")
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
				Id:          pItem.Id,
				OwnerId:     pItem.OwnerId,
				IsAutoTrack: pItem.IsAutoTrack,
				CategoryId:  pItem.CategoryId,
				Name:        pItem.Name,
				Image:       pItem.Image,
				Desc:        pItem.Desc,
				CreatedAt:   timestamppb.New(pItem.CreatedAt),
				DeletedAt:   conv.NullableTime(pItem.DeletedAt),
				UpdatedAt:   timestamppb.New(pItem.UpdatedAt),
				LastTracks:  []*pb.DBTrack{},
			}
			if pItem.CategoryId != nil {
				p.CategoryName = &pItem.Category.Name
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

func (d *DBCDeliveryService) TrackDay(ctx context.Context, r *pb.TrackDayRequest) (*pb.TrackDayResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "TrackDay")
	}

	date, err := tools.ParseISO(r.DateISO)
	if err != nil {
		return nil, errors.Wrap(err, "TrackDay")
	}

	uCaseRes, err := d.dbcChallengesUCase.TrackDay(ctx, &domain.DBCTrack{
		UserId:      userId,
		ChallengeId: r.ChallengeId,
		Date:        date,
		Done:        r.Done,
	})
	if err != nil {
		return nil, errors.Wrap(err, "TrackDay")
	}

	response := &pb.TrackDayResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.LastSeries = uCaseRes.LastSeries
		response.ScoreDaily = uCaseRes.ScoreDaily
	}

	return response, nil
}

func (d *DBCDeliveryService) GetMonthTracks(ctx context.Context, r *pb.GetMonthTracksRequest) (*pb.GetMonthTracksResponse, error) {
	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "ExtractRequestUserId")
	}

	date, err := tools.ParseISO(r.DateISO)
	if err != nil {
		return nil, errors.Wrap(err, "ParseISO")
	}

	uCaseRes, err := d.dbcChallengesUCase.GetMonthTracks(ctx, date, r.ChallengeId, userId)
	if err != nil {
		return nil, errors.Wrap(err, "GetMonthTracks")
	}

	response := &pb.GetMonthTracksResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.Tracks = []*pb.DBTrack{}
		for _, pTrack := range uCaseRes.Tracks {
			t := &pb.DBTrack{
				Date:       timestamppb.New(pTrack.Date),
				DateString: pTrack.Date.Format("02-01-2006"),
				Done:       pTrack.Done,
				LastSeries: pTrack.LastSeries,
				Score:      pTrack.Score,
				ScoreDaily: pTrack.ScoreDaily,
			}
			response.Tracks = append(response.Tracks, t)
		}
	}

	return response, nil
}

func (d *DBCDeliveryService) GetChallengeInfo(ctx context.Context, r *pb.IdRequest) (*pb.GetChallengeInfoResponse, error) {

	userId, err := app.ExtractRequestUserId(ctx)
	if err != nil {
		// Let's suppose that it is just public call [ But no way :) ]
		userId = -1
	}

	uCaseRes, err := d.dbcChallengesUCase.Info(userId, r.Id)
	if err != nil {
		return nil, errors.Wrap(err, "Info")
	}

	response := &pb.GetChallengeInfoResponse{
		Status: &pb.Status{
			Code:    uCaseRes.StatusCode,
			Message: uCaseRes.StatusCode,
		},
	}

	if uCaseRes.StatusCode == domain.Success {
		response.IsMember = uCaseRes.IsMember
		response.Challenge = &pb.DBCChallenge{
			Id:          uCaseRes.Challenge.Id,
			OwnerId:     uCaseRes.Challenge.OwnerId,
			CategoryId:  uCaseRes.Challenge.CategoryId,
			IsAutoTrack: uCaseRes.Challenge.IsAutoTrack,
			Name:        uCaseRes.Challenge.Name,
			Desc:        uCaseRes.Challenge.Desc,
			Image:       uCaseRes.Challenge.Image,
			CreatedAt:   timestamppb.New(uCaseRes.Challenge.CreatedAt),
			UpdatedAt:   timestamppb.New(uCaseRes.Challenge.UpdatedAt),
		}

		if uCaseRes.Challenge.Category != nil {
			response.Challenge.CategoryName = &uCaseRes.Challenge.Category.Name
		}
		if uCaseRes.Challenge.DeletedAt != nil {
			response.Challenge.DeletedAt = timestamppb.New(*uCaseRes.Challenge.DeletedAt)
		}
	}

	return response, nil
}
