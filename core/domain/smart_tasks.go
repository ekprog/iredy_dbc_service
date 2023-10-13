package domain

import "time"

type GenerationPeriod string

const (
	GenerationPeriodDay   GenerationPeriod = "day"
	GenerationPeriodWeek  GenerationPeriod = "week"
	GenerationPeriodMonth GenerationPeriod = "month"
)

type GenerationItem struct {
	Id       int32            `json:"id"`
	Period   GenerationPeriod `json:"period"`
	Datetime time.Time        `json:"datetime"`

	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type SmartTask struct {
	Id        int32  `json:"id"`
	UserId    int32  `json:"user_id"`
	ProjectId *int32 `json:"project_id"`

	Name     string   `json:"name"`
	Desc     *string  `json:"desc"`
	Priority Priority `json:"priority"`

	GenerationItems []*GenerationItem

	LastGeneratedAt *time.Time `json:"last_generated_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type SmartTasksRepository interface {

	// System
	CountAllActive() (int32, error)
	FetchAllActive(offset, limit int32) ([]*SmartTask, error)
	UpdateLastGeneratedAt(id int32, time time.Time) error

	// User
	FetchAllByUserId(userId int32, trashed bool, offset, limit int32) ([]*SmartTask, error)
}

type PeriodMatcher interface {
	Match(timeX time.Time, timeY time.Time, by GenerationPeriod) (bool, error)
}

type SmartTasksInteractor interface {
	All(userId int32, trashed bool, offset, limit int32) (SmartTasksListResponse, error)

	// Generate tasks for all users using smart task template
	GenerateTasks() (StatusResponse, error)
}

type SmartTasksListResponse struct {
	StatusCode string
	Tasks      []*SmartTask
}
