package domain

import "time"

type Project struct {
	Id     int32 `json:"id"`
	UserId int32 `json:"user_id"`

	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Color string `json:"color"`

	Order    int32  `json:"order"`
	ParentId *int32 `json:"parent_id"`

	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ProjectDrag struct {
	ProjectId int32
	ParentId  *int32
	Order     int32
}

type ProjectsRepository interface {
	FetchByUserId(int32) ([]*Project, error)
	FetchByUserIdTrashed(int32) ([]*Project, error)
	FetchById(int32) (*Project, error)

	Insert(*Project) error
	Update(*Project) error
	UpdateOrderForUser(int32, []int32) error
	DragItemsForUser(int32, []*ProjectDrag) error
	Remove(int32) error
}

type ProjectsInteractor interface {
	Active(userId int32, trashed bool) (ProjectListResponse, error)
	Info(userId int32, projectId int32) (ProjectInfoResponse, error)
	Trashed(userId int32) (ProjectListResponse, error)
	Create(project Project) (IdResponse, error)
	Remove(userId, projectId int32) (StatusResponse, error)
	Update(*Project) (StatusResponse, error)
	UpdateOrder(userId int32, newOrder []int32) (StatusResponse, error)
	DragProjects(userId int32, drags []*ProjectDrag) (StatusResponse, error)
}

// Responses (only for UseCase layer)

type ProjectListResponse struct {
	StatusCode string
	Projects   []*Project
}

type ProjectInfoResponse struct {
	StatusCode string
	Project    *Project
}
