package domain

type UserGamify struct {
	Score      int64
	ScoreDaily int64
}

// IO FORMS (RESPONSES)

type UserGamifyResponse struct {
	StatusCode string
	LastSeries int64
	ScoreDaily int64
}
