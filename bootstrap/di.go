package bootstrap

import (
	"go.uber.org/dig"
	"microservice/app"
	"microservice/core/delivery"
	"microservice/core/domain"
	"microservice/core/interactors"
	"microservice/core/repos"
	"microservice/core/services"
)

func initDependencies(di *dig.Container) error {

	// Repository
	di.Provide(repos.NewUsersRepo, dig.As(new(domain.UsersRepository)))
	di.Provide(repos.NewProjectsRepo, dig.As(new(domain.ProjectsRepository)))
	di.Provide(repos.NewTasksRepo, dig.As(new(domain.TasksRepository)))
	di.Provide(repos.NewSmartTasksRepo, dig.As(new(domain.SmartTasksRepository)))

	// Services
	di.Provide(services.NewPeriodMatcherService, dig.As(new(domain.PeriodMatcher)))

	// Use Cases
	di.Provide(interactors.NewUsersInteractor, dig.As(new(domain.UsersInteractor)))
	di.Provide(interactors.NewProjectsInteractor, dig.As(new(domain.ProjectsInteractor)))
	di.Provide(interactors.NewTasksInteractor, dig.As(new(domain.TasksInteractor)))
	di.Provide(interactors.NewSmartTasksInteractor, dig.As(new(domain.SmartTasksInteractor)))

	di.Provide(delivery.NewStatusDeliveryService)
	di.Provide(delivery.NewToDoDeliveryService)

	// Delivery
	if err := app.InitDelivery(delivery.NewStatusDeliveryService); err != nil {
		return err
	}

	if err := app.InitDelivery(delivery.NewToDoDeliveryService); err != nil {
		return err
	}

	return nil
}
