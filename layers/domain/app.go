package domain

// RESPONSE CODES
const (
	Success         string = "success"
	ValidationError string = "validation_error"
	NotFound        string = "not_found"
	AlreadyExists   string = "already_exists"
	UserLogicError         = "user_error"
	ServerError            = "server_error"
)

// GENERAL RESPONSES
type StatusResponse struct {
	StatusCode string
}

type IdResponse struct {
	StatusCode string
	Id         int64
}

// SCORES
type ScorePoints struct {
	Score      int64
	ScoreDaily int64
}

// PERIOD_TYPE
type PeriodType = string

const (
	PeriodTypeEveryDay   PeriodType = "every_day"
	PeriodTypeWeekDays   PeriodType = "week_days"
	PeriodTypeMonthDates PeriodType = "month_dates"
)

type GenerationPeriod struct {
	Type PeriodType
	Data []int
}
