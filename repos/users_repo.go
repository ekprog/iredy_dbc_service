package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/domain"
)

type UsersRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewUsersRepo(log core.Logger, db *sql.DB) *UsersRepo {
	return &UsersRepo{log: log, db: db}
}

func (r *UsersRepo) Exist(id int32) (bool, error) {
	query := `select id from users where id=$1 limit 1`
	err := r.db.QueryRow(query, id).Scan(&id)
	switch err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (r *UsersRepo) InsertIfNotExists(user *domain.User) error {
	query := `INSERT INTO users (id) VALUES ($1) ON CONFLICT DO NOTHING;`
	_, err := r.db.Exec(query, user.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UsersRepo) Remove(id int32) error {
	query := `UPDATE users set deleted_at=now() where id=$1;`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
