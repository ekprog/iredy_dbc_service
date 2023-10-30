package repos

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	"microservice/app/core"
	"microservice/layers/do"
	"microservice/layers/domain"
)

type DBCUserChallengesRepo struct {
	log    core.Logger
	db     *sql.DB
	gormDB *gorm.DB
}

func NewDBCUserChallengesRepo(log core.Logger, db *sql.DB, gormDB *gorm.DB) *DBCUserChallengesRepo {
	return &DBCUserChallengesRepo{log: log, db: db, gormDB: gormDB}
}

func (r *DBCUserChallengesRepo) FetchAll(limit, offset int64) ([]*domain.DBCUserChallenge, error) {

	var items []*do.DBCChallengesUsers
	err := r.gormDB.Table("dbc_challenges_users").
		Preload("User").
		Preload("Challenge.Category").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCUserChallenge
	for _, item := range items {
		result = append(result, item.DTO())
	}

	return result, nil
}

func (r *DBCUserChallengesRepo) UserFetchAll(userId int64) ([]*domain.DBCUserChallenge, error) {

	var items []*do.DBCChallengesUsers
	err := r.gormDB.
		Table("dbc_challenges_users").
		Preload("Challenge").
		Preload("Challenge.Category").
		Where(&do.DBCChallengeCategory{UserID: userId}).
		Order("created_at desc").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCUserChallenge
	for _, item := range items {
		result = append(result, item.DTO())
	}

	return result, nil
}

func (r *DBCUserChallengesRepo) FetchById(ctx context.Context, id int64) (*domain.DBCUserChallenge, error) {
	var item = &domain.DBCUserChallenge{
		Id: id,
		ChallengeInfo: &domain.DBCChallengeInfo{
			Category: &domain.DBCCategory{},
		},
	}
	query := `select 
    			c.id,
    			c.user_id,
    			c.challenge_id,
    			ci.is_auto_track,
    			ci."desc", 
    			c.created_at, 
    			c.updated_at,
    			c.deleted_at
			from dbc_challenges_users c
					left join dbc_challenges ci on ci.id = c.challenge_id
				where c.id=$1
			limit 1`

	var categoryId *int64
	var categoryName *string

	err := r.db.QueryRow(query, id).Scan(
		&item.Id,
		&item.UserId,
		&item.ChallengeInfoId,
		&item.ChallengeInfo.IsAutoTrack,
		&item.ChallengeInfo.Desc,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt)

	if err == nil && categoryId != nil && categoryName != nil {
		category := &domain.DBCCategory{
			Id:   *categoryId,
			Name: *categoryName,
		}
		item.ChallengeInfo.Category = category
		item.ChallengeInfo.Id = item.ChallengeInfoId
	}

	switch err {
	case nil:
		return item, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
}

func (r *DBCUserChallengesRepo) UserFetchByName(userId int64, name string) (*domain.DBCUserChallenge, error) {
	var item = &domain.DBCUserChallenge{
		UserId: userId,
		ChallengeInfo: &domain.DBCChallengeInfo{
			Name:     name,
			Category: &domain.DBCCategory{},
		},
	}
	query := `select 
    			c.id,
    			ci.category_id, 
    			cat.name as category_name,
    			ci.is_auto_track,
    			ci."desc", 
    			c.created_at, 
    			c.updated_at,
    			c.deleted_at
			from dbc_challenges_users c
					left join dbc_challenges ci on ci.id = c.challenge_id
					left join dbc_challenge_categories cat on ci.category_id = cat.id
			where c.user_id=$1 and ci.name=$2
			limit 1`
	err := r.db.QueryRow(query, userId, name).Scan(
		&item.Id,
		&item.ChallengeInfo.Category.Id,
		&item.ChallengeInfo.Category.Name,
		&item.ChallengeInfo.IsAutoTrack,
		&item.ChallengeInfo.Desc,
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

func (r *DBCUserChallengesRepo) Insert(item *domain.DBCUserChallenge) error {

	doItem := &do.DBCChallengesUsers{
		UserID:      item.UserId,
		ChallengeID: item.ChallengeInfo.Id,
		LastSeries:  0,
	}
	err := r.gormDB.Table("dbc_challenges_users").Create(doItem).Error
	if err != nil {
		return err
	}

	item.Id = int64(doItem.ID)
	return nil
}

func (r *DBCUserChallengesRepo) Update(item *domain.DBCUserChallenge) error {
	query := `UPDATE dbc_challenges_users 
				SET last_series=$2, updated_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query,
		item.Id,
		item.LastSeries)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCUserChallengesRepo) Remove(id int64) error {
	query := `UPDATE dbc_challenges 
				SET deleted_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCUserChallengesRepo) UserExistsByChallengeId(userId, challengeId int64) (bool, error) {
	query := `select count(id) from dbc_challenges_users 
                 where user_id = $1 and challenge_id = $2`

	var c int64
	err := r.db.QueryRow(query, userId, challengeId).Scan(&c)
	if err != nil {
		return false, err
	}

	return c > 0, nil
}
