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

//func (r *DBCChallengesRepo) FetchAll(userId int32) ([]*domain.DBCChallenge, error) {
//
//	query := `select id,
//					 category_id,
//					 name,
//					 "desc",
//					 image,
//					 last_series,
//					 created_at,
//					 updated_at,
//					 deleted_at,
//					 date,
//					 done
//				from (select c.id,
//							 c.category_id,
//							 c.name,
//							 c."desc",
//							 c.image,
//							 c.last_series,
//							 c.created_at,
//							 c.updated_at,
//							 c.deleted_at,
//							 t.date,
//							 t.done,
//							 row_number() over (PARTITION BY c.id ORDER BY t.date) as num
//					  from dbc_challenges c
//							   left join dbc_challenges_tracks t
//										 on c.id = t.challenge_id
//					  where c.user_id = $1
//						and c.deleted_at is null
//					  order by c.id, t.date DESC) all_t
//				where all_t.num <= 3`
//
//	rows, err := r.db.Query(query, userId)
//	if err != nil {
//		return nil, err
//	}
//
//	result := make(map[int32]*domain.DBCChallenge)
//	var order []int32
//
//	for rows.Next() {
//		var date *time.Time
//		var done *bool
//
//		item := &domain.DBCChallenge{
//			UserId:     userId,
//			LastTracks: []*domain.DBCTrack{},
//		}
//		err := rows.Scan(&item.Id,
//			&item.CategoryId,
//			&item.Name,
//			&item.Desc,
//			&item.Image,
//			&item.LastSeries,
//			&item.CreatedAt,
//			&item.UpdatedAt,
//			&item.DeletedAt,
//			&date,
//			&done)
//		if err != nil {
//			return nil, err
//		}
//
//		if _, ok := result[item.Id]; !ok {
//			result[item.Id] = item
//			order = append(order, item.Id)
//		}
//
//		if date != nil && done != nil {
//			result[item.Id].LastTracks = append(result[item.Id].LastTracks, &domain.DBCTrack{
//				Date: *date,
//				Done: *done,
//			})
//		}
//	}
//
//	// To array
//	values := make([]*domain.DBCChallenge, 0, len(result))
//	for _, ind := range order {
//		values = append(values, result[ind])
//	}
//
//	return values, nil
//}

func (r *DBCChallengesRepo) FetchById(id int32) (*domain.DBCChallenge, error) {
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
	err := r.db.QueryRow(query, id).Scan(
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

func (r *DBCChallengesRepo) FetchByName(userId int32, name string) (*domain.DBCChallenge, error) {
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
