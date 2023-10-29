package repos

import (
	"database/sql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"microservice/app/core"
	"microservice/layers/do"
	"microservice/layers/domain"
)

type DBCChallengesRepo struct {
	log    core.Logger
	db     *sql.DB
	gormDB *gorm.DB
}

func NewDBCChallengesRepo(log core.Logger, db *sql.DB, gormDB *gorm.DB) *DBCChallengesRepo {
	return &DBCChallengesRepo{log: log, db: db, gormDB: gormDB}
}

func (r *DBCChallengesRepo) Insert(item *domain.DBCChallengeInfo) error {

	doItem, err := do.NewDBCChallenge(item)
	if err != nil {
		return errors.Wrap(err, "NewDBCChallenge")
	}

	err = r.gormDB.Table("dbc_challenges").Create(doItem).Error
	if err != nil {
		return err
	}

	item.Id = int64(doItem.ID)
	return nil
}
