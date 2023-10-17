package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/layers/domain"
	"time"
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
    			c.id,
    			c.category_id,
    			c.name, 
    			c."desc", 
    			c.last_series,
    			c.created_at, 
    			c.updated_at,
    			c.deleted_at,
    			t.date
			from dbc_challenges c
				left join dbc_challenges_tracks t 
				    on c.id = t.challenge_id
			where c.user_id=$1 and 
			      c.deleted_at is null and 
			      t.deleted_at is null and
				  (
				      t.id is null or
				      t.id IN (SELECT ct.id
                 		FROM dbc_challenges_tracks ct
                 		order by ct.date desc
                 		LIMIT 5 OFFSET 0)
				 )
			order by c.created_at desc, t.date desc`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	result := make(map[int32]*domain.DBCChallenge)
	var order []int32

	for rows.Next() {
		var date *time.Time
		item := &domain.DBCChallenge{
			UserId:     userId,
			LastTracks: []*domain.DBCTrack{},
		}
		err := rows.Scan(&item.Id,
			&item.CategoryId,
			&item.Name,
			&item.Desc,
			&item.LastSeries,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
			&date)
		if err != nil {
			return nil, err
		}

		if _, ok := result[item.Id]; !ok {
			result[item.Id] = item
			order = append(order, item.Id)
		}

		if date != nil {
			result[item.Id].LastTracks = append(result[item.Id].LastTracks, &domain.DBCTrack{Date: *date})
		}
	}

	// To array
	values := make([]*domain.DBCChallenge, 0, len(result))
	for _, ind := range order {
		values = append(values, result[ind])
	}

	return values, nil
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

func (r *DBCChallengesRepo) FetchByName(userId int32, name string) (*domain.DBCChallenge, error) {
	var task = &domain.DBCChallenge{
		UserId: userId,
		Name:   name,
	}
	query := `select 
    			id,
    			category_id, 
    			"desc", 
    			created_at, 
    			updated_at,
    			deleted_at
			from dbc_challenges
			where user_id=$1 and name=$2
			limit 1`
	err := r.db.QueryRow(query, userId, name).Scan(
		&task.Id,
		&task.CategoryId,
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
