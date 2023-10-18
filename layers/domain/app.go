package domain

// RESPONSE CODES
const (
	Success         string = "success"
	ValidationError string = "validation_error"
	NotFound        string = "not_found"
	AlreadyExists   string = "already_exists"
)

// GENERAL RESPONSES
type StatusResponse struct {
	StatusCode string
}

type IdResponse struct {
	StatusCode string
	Id         int32
}

// PERIOD_TYPE
type PeriodType = string

const (
	PeriodTypeEveryDay PeriodType = "every_day"
)
