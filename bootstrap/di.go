package bootstrap

import (
	"go.uber.org/dig"
	"microservice/app"
	"microservice/layers/delivery/grpc"
	"microservice/layers/domain"
	"microservice/layers/repos"
	"microservice/layers/services"
	"microservice/layers/usecase"
)

func initDependencies(di *dig.Container) error {

	// Repository
	_ = di.Provide(repos.NewUsersRepo, dig.As(new(domain.UsersRepository)))
	_ = di.Provide(repos.NewDBCTracksRepo, dig.As(new(domain.DBCTrackRepository)))
	_ = di.Provide(repos.NewDBCCategoriesRepo, dig.As(new(domain.DBCCategoryRepository)))
	_ = di.Provide(repos.NewDBCChallengesRepo, dig.As(new(domain.DBCChallengesRepository)))

	// Services
	_ = di.Provide(services.NewPeriodTypeProcessor)
	_ = di.Provide(services.NewDBCTrackProcessor)
	_ = di.Provide(services.NewAchievementsProcessor)

	// Use Cases
	_ = di.Provide(usecase.NewUsersUseCase, dig.As(new(domain.UsersUseCase)))
	_ = di.Provide(usecase.NewDBCCategoriesUCase, dig.As(new(domain.DBCCategoryUseCase)))
	_ = di.Provide(usecase.NewChallengesUseCase, dig.As(new(domain.DBCChallengesUseCase)))

	_ = di.Provide(grpc.NewStatusDeliveryService)
	_ = di.Provide(grpc.NewDBCDeliveryService)

	// Delivery
	if err := app.InitDelivery(grpc.NewStatusDeliveryService); err != nil {
		return err
	}

	if err := app.InitDelivery(grpc.NewUsersDeliveryService); err != nil {
		return err
	}

	if err := app.InitDelivery(grpc.NewDBCDeliveryService); err != nil {
		return err
	}

	return nil
}
