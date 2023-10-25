package repos

import (
	"context"
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

func (r *DBCChallengesRepo) FetchAll(limit, offset int64) ([]*domain.DBCChallenge, error) {

	query := `select c.id,
    			 c.user_id,
				 c.category_id,
				 cat.name as category_name,
				 c.is_auto_track,
				 c.name,
				 c."desc",
				 c.image,
				 c.last_series,
				 c.created_at,
				 c.updated_at,
				 c.deleted_at
		  from dbc_challenges c
		  		left join dbc_challenge_categories cat on c.category_id = cat.id
		  where c.deleted_at is null 
		  order by c.id
		  limit $1 offset $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCChallenge

	for rows.Next() {
		item := &domain.DBCChallenge{}
		err := rows.Scan(&item.Id,
			&item.UserId,
			&item.CategoryId,
			&item.CategoryName,
			&item.IsAutoTrack,
			&item.Name,
			&item.Desc,
			&item.Image,
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

func (r *DBCChallengesRepo) FetchUsersAll(userId int64) ([]*domain.DBCChallenge, error) {

	query := `select c.id,
				 c.category_id,
				 cat.name as category_name,
				 c.is_auto_track,
				 c.name,
				 c."desc",
				 c.image,
				 c.last_series,
				 c.created_at,
				 c.updated_at,
				 c.deleted_at
		  from dbc_challenges c
		  		left join dbc_challenge_categories cat on c.category_id = cat.id
		  where c.user_id = $1
			and c.deleted_at is null
		  order by c.id`

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
			&item.CategoryName,
			&item.IsAutoTrack,
			&item.Name,
			&item.Desc,
			&item.Image,
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

func (r *DBCChallengesRepo) FetchById(ctx context.Context, id int64) (*domain.DBCChallenge, error) {
	var item = &domain.DBCChallenge{
		Id: id,
	}
	query := `select 
    			c.name,
    			c.user_id,
    			c.category_id, 
    			cat.name as category_name,
    			c.is_auto_track,
    			c."desc", 
    			c.created_at, 
    			c.updated_at,
    			c.deleted_at
			from dbc_challenges c
					left join dbc_challenge_categories cat on c.category_id = cat.id
			where c.id=$1
			limit 1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.Name,
		&item.UserId,
		&item.CategoryId,
		&item.CategoryName,
		&item.IsAutoTrack,
		&item.Desc,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt)
	switch err {
	case nil:
		return item, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *DBCChallengesRepo) FetchByName(userId int64, name string) (*domain.DBCChallenge, error) {
	var item = &domain.DBCChallenge{
		UserId: userId,
		Name:   name,
	}
	query := `select 
    			c.id,
    			c.category_id, 
    			cat.name as category_name,
    			c.is_auto_track,
    			c."desc", 
    			c.created_at, 
    			c.updated_at,
    			c.deleted_at
			from dbc_challenges c
					left join dbc_challenge_categories cat on c.category_id = cat.id
			where c.user_id=$1 and c.name=$2
			limit 1`
	err := r.db.QueryRow(query, userId, name).Scan(
		&item.Id,
		&item.CategoryId,
		&item.CategoryName,
		&item.IsAutoTrack,
		&item.Desc,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt)
	switch err {
	case nil:
		return item, nil
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
                   is_auto_track,
                   name, 
                   "desc") 
			 VALUES ($1, $2, $3, $4, $5) returning id;`
	err := r.db.QueryRow(query,
		item.UserId,
		item.CategoryId,
		item.IsAutoTrack,
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
				SET name=$2, "desc"=$3, last_series=$4, updated_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query,
		item.Id,
		item.Name,
		item.Desc,
		item.LastSeries)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCChallengesRepo) Remove(id int64) error {
	query := `UPDATE dbc_challenges 
				SET deleted_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
