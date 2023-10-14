package bootstrap

import (
	"go.uber.org/dig"
	"microservice/app"
	"microservice/layers/delivery/grpc"
	"microservice/layers/domain"
	"microservice/layers/repos"
	"microservice/layers/usecase"
)

func initDependencies(di *dig.Container) error {

	// Repository
	di.Provide(repos.NewUsersRepo, dig.As(new(domain.UsersRepository)))
	di.Provide(repos.NewDBCCategoriesRepo, dig.As(new(domain.DBCCategoryRepository)))
	di.Provide(repos.NewDBCChallengesRepo, dig.As(new(domain.DBCChallengesRepository)))

	// Use Cases
	di.Provide(usecase.NewUsersUseCase, dig.As(new(domain.UsersUseCase)))
	di.Provide(usecase.NewDBCCategoriesUCase, dig.As(new(domain.DBCCategoryUseCase)))
	di.Provide(usecase.NewChallengesUseCase, dig.As(new(domain.ChallengesUseCase)))

	di.Provide(grpc.NewStatusDeliveryService)
	di.Provide(grpc.NewDBCDeliveryService)

	// Delivery
	if err := app.InitDelivery(grpc.NewStatusDeliveryService); err != nil {
		return err
	}

	if err := app.InitDelivery(grpc.NewDBCDeliveryService); err != nil {
		return err
	}

	return nil
}
