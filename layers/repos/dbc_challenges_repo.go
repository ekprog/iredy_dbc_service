package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/layers/domain"
)

type DBCChallengesRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewDBCChallengesRepo(log core.Logger, db *sql.DB) *DBCChallengesRepo {
	return &DBCChallengesRepo{log: log, db: db}
}

func (r *DBCChallengesRepo) FetchAll(userId int32) ([]*domain.DBCChallenge, error) {

	query := `select 
    			id,
    			category_id,
    			name, 
    			"desc", 
    			last_series,
    			created_at, 
    			updated_at,
    			deleted_at
			from dbc_challenges
			where user_id=$1 and deleted_at is null
			order by created_at`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCChallenge
	for rows.Next() {
		item := &domain.DBCChallenge{
			UserId: userId,
		}
		err := rows.Scan(&item.Id,
			&item.CategoryId,
			&item.Name,
			&item.Desc,
			&item.LastSeries,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *DBCChallengesRepo) FetchById(id int32) (*domain.DBCChallenge, error) {
	var task = &domain.DBCChallenge{
		Id: id,
	}
	query := `select 
    			user_id, 
    			project_id, 
    			name, 
    			"desc", 
    			priority, 
    			done,
    			created_at, 
    			updated_at,
    			deleted_at
			from tasks
			where id=$1
			limit 1`
	err := r.db.QueryRow(query, id).Scan(&task.UserId,
		&task.CategoryId,
		&task.Name,
		&task.Desc,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt)
	switch err {
	case nil:
		return task, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *DBCChallengesRepo) Insert(item *domain.DBCChallenge) error {
	query := `INSERT INTO dbc_challenges (
                   user_id, 
                   category_id, 
                   name, 
                   "desc") 
			 VALUES ($1, $2, $3, $4) returning id;`
	err := r.db.QueryRow(query,
		item.UserId,
		item.CategoryId,
		item.Name,
		item.Desc,
	).Scan(&item.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCChallengesRepo) Update(item *domain.DBCChallenge) error {
	query := `UPDATE dbc_challenges 
				SET name=$2, "desc"=$3, updated_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query,
		item.Id,
		item.Name,
		item.Desc)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCChallengesRepo) Remove(id int32) error {
	query := `UPDATE dbc_challenges 
				SET deleted_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
