package domain

import "time"

// ToDo: DTO
type Challenge struct {
	Id         int32  `json:"id"`
	UserId     int32  `json:"user_id"`
	CategoryId *int32 `json:"category_id"`

	Name       string `json:"name"`
	Desc       string `json:"desc"`
	LastSeries int32

	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CreateChallengeForm struct {
	Id           int32  `json:"id"`
	UserId       int32  `json:"user_id"`
	Name         string `json:"name"`
	CategoryName *string
}

type ChallengesRepository interface {
	FetchAll(userId int32) ([]*Challenge, error)
	FetchById(int32) (*Challenge, error)
	Insert(*Challenge) error
	Update(*Challenge) error
	Remove(int32) error
}

type ChallengesInteractor interface {
	All(userId int32) (ChallengesListResponse, error)
	Create(form *CreateChallengeForm) (IdResponse, error)
	Update(task *Challenge) (StatusResponse, error)
	Remove(userId, taskId int32) (StatusResponse, error)
}

type ChallengesListResponse struct {
	StatusCode string
	Challenges []*Challenge
}
