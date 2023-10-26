package services

import (
	"context"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"microservice/app/core"
	"microservice/layers/domain"
)

type AchievementsProcessor struct {
	log        core.Logger
	trxManager *manager.Manager
	userRepo   domain.UsersRepository
}

func NewAchievementsProcessor(log core.Logger,
	trxManager *manager.Manager,
	userRepo domain.UsersRepository) *AchievementsProcessor {
	return &AchievementsProcessor{
		log:        log,
		trxManager: trxManager,
		userRepo:   userRepo,
	}
}

// When user`s score changed
func (s *AchievementsProcessor) HandleUserScoreChanged(ctx context.Context, userId, score int64) error {
	// Here should be triggers that start executing after some achievement point

	return nil
}

func (s *AchievementsProcessor) HandleChallengeLastSeriesChanged(ctx context.Context, challengeId, lastSeries int64) error {
	// Here should be triggers that start executing after some achievement point

	return nil
}
