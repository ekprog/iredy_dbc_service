package domain

type UserGamify struct {
	Score      int32
	ScoreDaily int32
}

// IO FORMS (RESPONSES)

type UserGamifyResponse struct {
	StatusCode string
	LastSeries int32
	ScoreDaily int32
}
