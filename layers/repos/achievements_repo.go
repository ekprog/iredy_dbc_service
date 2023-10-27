package repos

import (
	"database/sql"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"microservice/app/core"
)

type AchievementsRepo struct {
	log    core.Logger
	db     *sql.DB
	getter *trmsql.CtxGetter
}

func NewAchievementsRepo(log core.Logger, db *sql.DB, getter *trmsql.CtxGetter) *AchievementsRepo {
	return &AchievementsRepo{log: log, db: db, getter: getter}
}
