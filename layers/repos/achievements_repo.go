package repos

import (
	"database/sql"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"github.com/pkg/errors"
	"microservice/app/core"
	"microservice/layers/domain"
)

type AchievementsRepo struct {
	log    core.Logger
	db     *sql.DB
	getter *trmsql.CtxGetter
}

func NewAchievementsRepo(log core.Logger, db *sql.DB, getter *trmsql.CtxGetter) *AchievementsRepo {
	return &AchievementsRepo{log: log, db: db, getter: getter}
}

func (r *AchievementsRepo) FetchById(id int64) (*domain.Achievement, error) {
	query := `select 
    				score,
    				created_at, 
    				updated_at, 
    				deleted_at
				from users where id=$1 limit 1`

	user := &domain.User{Id: id}

	err := r.db.QueryRow(query, id).Scan(
		&user.Score,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "FetchById")
	}
}
