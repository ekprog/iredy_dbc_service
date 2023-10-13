package domain

import "time"

type Priority int

const (
	Priority1 Priority = 0
	Priority2 Priority = 1
	Priority3 Priority = 2
	Priority4 Priority = 3
)

type Task struct {
	Id        int32  `json:"id"`
	UserId    int32  `json:"user_id"`
	ProjectId *int32 `json:"project_id"`

	Name     string   `json:"name"`
	Desc     *string  `json:"desc"`
	Priority Priority `json:"priority"`
	Done     bool     `json:"done"`

	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type TasksRepository interface {
	FetchAll(userId int32, done bool, offset, limit int32) ([]*Task, error)
	FetchById(int32) (*Task, error)
	FetchByProjectId(projectId int32, done bool, offset, limit int32) ([]*Task, error)
	Insert(*Task) error
	Update(*Task) error
	Remove(int32) error
}

type TasksInteractor interface {
	All(userId int32, projectId *int32, done bool, offset, limit int32) (TasksListResponse, error)
	Create(task *Task) (IdResponse, error)
	Update(task *Task) (StatusResponse, error)
	SetDone(userId, taskId int32, flag bool) (StatusResponse, error)
	Remove(userId, taskId int32) (StatusResponse, error)
}

type TasksListResponse struct {
	StatusCode string
	Tasks      []*Task
}
