package services

import (
	"context"
	"microservice/app/core"
	"time"
)

type DBCTrackProcessor struct {
	log            core.Logger
	trackProcessor *PeriodTypeProcessor
}

func NewDBCTrackProcessor(log core.Logger, trackProcessor *PeriodTypeProcessor) *DBCTrackProcessor {
	return &DBCTrackProcessor{
		log:            log,
		trackProcessor: trackProcessor,
	}
}

// Меняет значение трека и всей предыдущей цепочки треков
func (s *DBCTrackProcessor) MakeTrack(ctx context.Context, challengeId int64, date time.Time, value bool) error {
	date = date.UTC()

	return nil
}
