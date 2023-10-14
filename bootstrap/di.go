package bootstrap

import (
	"go.uber.org/dig"
	"microservice/app"
	"microservice/layers/delivery"
	"microservice/layers/domain"
	"microservice/layers/interactors"
	"microservice/layers/repos"
)

func initDependencies(di *dig.Container) error {

	// Repository
	di.Provide(repos.NewUsersRepo, dig.As(new(domain.UsersRepository)))
	di.Provide(repos.NewDBCCategoriesRepo, dig.As(new(domain.DBCCategoryRepository)))
	di.Provide(repos.NewDBCChallengesRepo, dig.As(new(domain.DBCChallengesRepository)))

	// Use Cases
	di.Provide(interactors.NewUsersInteractor, dig.As(new(domain.UsersInteractor)))
	di.Provide(interactors.NewDBCCategoriesUCase, dig.As(new(domain.DBCCategoryInteractor)))
	di.Provide(interactors.NewChallengesInteractor, dig.As(new(domain.ChallengesInteractor)))

	di.Provide(delivery.NewStatusDeliveryService)
	di.Provide(delivery.NewDBCDeliveryService)

	// Delivery
	if err := app.InitDelivery(delivery.NewStatusDeliveryService); err != nil {
		return err
	}

	if err := app.InitDelivery(delivery.NewDBCDeliveryService); err != nil {
		return err
	}

	return nil
}
