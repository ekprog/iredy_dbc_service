package domain

import (
	"time"
)

//
// MODELS
//

type DBCCategory struct {
	Id     int32
	UserId int32

	Name string

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCChallenge struct {
	Id         int32
	UserId     int32
	CategoryId *int32

	Name       string
	Desc       string
	LastSeries int32

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

//
// REPOSITORIES
//

type DBCCategoryRepository interface {
	FetchByUserId(int32) ([]*DBCCategory, error)
	FetchByName(int32, string) (*DBCCategory, error)
	FetchById(int32) (*DBCCategory, error)
	Insert(*DBCCategory) error
	Update(*DBCCategory) error
	Remove(int) error
}

type DBCChallengesRepository interface {
	FetchAll(userId int32) ([]*DBCChallenge, error)
	FetchById(int32) (*DBCChallenge, error)
	FetchByName(int32, string) (*DBCChallenge, error)
	Insert(*DBCChallenge) error
	Update(*DBCChallenge) error
	Remove(int32) error
}

//
// USE CASES
//

type DBCCategoryUseCase interface {
	Get(userId int32) (CategoryListResponse, error)
	Update(*DBCCategory) (StatusResponse, error)
}

type ChallengesUseCase interface {
	All(userId int32) (ChallengesListResponse, error)
	Create(form *CreateDBCChallengeForm) (IdResponse, error)
	Update(task *DBCChallenge) (StatusResponse, error)
	Remove(userId, taskId int32) (StatusResponse, error)
}

// IO FORMS

type CreateDBCChallengeForm struct {
	UserId       int32
	Name         string
	Desc         string
	CategoryName *string
}

type ChallengesListResponse struct {
	StatusCode string
	Challenges []*DBCChallenge
}

type CategoryListResponse struct {
	StatusCode string
	Categories []*DBCCategory
}
