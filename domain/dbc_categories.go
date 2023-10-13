package domain

import "time"

type DBCCategory struct {
	Id     int32 `json:"id"`
	UserId int32 `json:"user_id"`

	Name string `json:"name"`

	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type DBCCategoryRepository interface {
	FetchByUserId(int32) ([]*DBCCategory, error)
	FetchById(int32) (*DBCCategory, error)
	Insert(*DBCCategory) error
	Update(*DBCCategory) error
	Remove(int) error
}

type DBCCategoryInteractor interface {
	Get(userId int32) (CategoryListResponse, error)
	Update(*DBCCategory) (StatusResponse, error)
}

// Responses (only for UseCase layer)

type CategoryListResponse struct {
	StatusCode string
	Categories []*DBCCategory
}
