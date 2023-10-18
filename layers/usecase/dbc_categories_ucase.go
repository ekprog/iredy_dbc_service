package usecase

import (
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
)

type DBCCategoriesUCase struct {
	log            core.Logger
	categoriesRepo domain.DBCCategoryRepository
	usersRepo      domain.UsersRepository
	tasksRepo      domain.DBCChallengesRepository
}

func NewDBCCategoriesUCase(log core.Logger,
	usersRepo domain.UsersRepository,
	categoriesRepo domain.DBCCategoryRepository) *DBCCategoriesUCase {
	return &DBCCategoriesUCase{
		log:            log,
		usersRepo:      usersRepo,
		categoriesRepo: categoriesRepo,
	}
}

func (i *DBCCategoriesUCase) Get(userId int32) (domain.CategoryListResponse, error) {

	var categories []*domain.DBCCategory
	var err error

	categories, err = i.categoriesRepo.FetchByUserId(userId)
	if err != nil {
		return domain.CategoryListResponse{}, errors.Wrap(err, "cannot fetch categories by user id")
	}

	return domain.CategoryListResponse{
		StatusCode: domain.Success,
		Categories: categories,
	}, nil
}

func (i *DBCCategoriesUCase) Update(item *domain.DBCCategory) (domain.StatusResponse, error) {

	// Check if user's project exists
	project, err := i.categoriesRepo.FetchById(item.Id)
	if err != nil {
		return domain.StatusResponse{},
			errors.Wrapf(err, "cannot fetch category before updating. CategoryId=%d", item.Id)
	}

	if project == nil || project.UserId != item.UserId {
		return domain.StatusResponse{
			StatusCode: domain.NotFound,
		}, nil
	}

	// Update
	err = i.categoriesRepo.Update(item)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrapf(err, "cannot update category %d", project.Id)
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}

func (i *DBCCategoriesUCase) Remove(userId, Id int32) (domain.StatusResponse, error) {

	var err error

	err = i.categoriesRepo.Remove(userId, Id)
	if err != nil {
		return domain.StatusResponse{}, errors.Wrap(err, "cannot delete categories by user id")
	}

	return domain.StatusResponse{
		StatusCode: domain.Success,
	}, nil
}
