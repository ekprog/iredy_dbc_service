package domain

const (
	Success         string = "success"
	AccessDenied    string = "access_denied"
	ValidationError string = "validation_error"
	NotFound        string = "not_found"

	ProjectNotFound string = "project_not_found"
	TaskNotFound    string = "task_not_found"
)

type StatusResponse struct {
	StatusCode string
}

type IdResponse struct {
	StatusCode string
	Id         int32
}
