package repos

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"microservice/app/core"
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

func (r *DBCChallengesRepo) PublicFetchLike(search string, categoryId *int64, limit, offset int64) ([]*domain.DBCChallengeInfo, error) {

	whereClause := fmt.Sprintf("visibility_type = 'public' and c.name LIKE '%%%s%%'", search)
	if categoryId != nil {
		whereClause += fmt.Sprintf(" and c.category_id=%d", *categoryId)
	}

	query := fmt.Sprintf(`select 
    				c.id,
    				c.owner_id,
      				c.category_id,
      				dcc.name,
					c.is_auto_track,
					c.name,
					c.image,
					c."desc",
    				c.created_at, 
    				c.updated_at, 
    				c.deleted_at
				from dbc_challenges c
					left join dbc_challenge_categories dcc on c.category_id = dcc.id
				where %s
				limit $1 offset $2`, whereClause)

	var items []*domain.DBCChallengeInfo

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		item := &domain.DBCChallengeInfo{
			VisibilityType: "public",
		}

		var categoryName *string
		err := rows.Scan(
			&item.Id,
			&item.OwnerId,
			&item.CategoryId,
			&categoryName,
			&item.IsAutoTrack,
			&item.Name,
			&item.Image,
			&item.Desc,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		if categoryName != nil {
			item.Category = &domain.DBCCategory{
				Id:   *categoryId,
				Name: *categoryName,
			}
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *DBCChallengesRepo) Insert(item *domain.DBCChallengeInfo) error {

	query := `INSERT INTO dbc_challenges (
                            owner_id, 
                            category_id,
                            name, 
                            "desc",
                            is_auto_track,
                            visibility_type) VALUES ($1, $2, $3, $4, $5, $6) 
                                             RETURNING id`
	err := r.db.QueryRow(query,
		item.OwnerId,
		item.CategoryId,
		item.Name,
		item.Desc,
		item.IsAutoTrack,
		item.VisibilityType).Scan(&item.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCChallengesRepo) FetchById(id int64) (*domain.DBCChallengeInfo, error) {
	query := `select 
    	c.id,
    	c.name,
    	c.image,
    	c.desc,
    	c.category_id,
		dcc.name,
		c.visibility_type,
		c.is_auto_track,
		c.owner_id,
		c.created_at,
		c.updated_at,
		c.deleted_at
    	from dbc_challenges c 
    		left join dbc_challenge_categories dcc on c.category_id = dcc.id
		where c.id = $1 and c.deleted_at is null`

	item := &domain.DBCChallengeInfo{
		Id: id,
	}

	var categoryName *string
	err := r.db.QueryRow(query, id).Scan(
		&item.Id,
		&item.Name,
		&item.Image,
		&item.Desc,
		&item.CategoryId,
		&categoryName,
		&item.VisibilityType,
		&item.IsAutoTrack,
		&item.OwnerId,
		&item.CreatedAt,
		&item.UpdatedAt,
		&item.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	if categoryName != nil && item.CategoryId != nil {
		item.Category = &domain.DBCCategory{
			Id:   *item.CategoryId,
			Name: *categoryName,
		}
	}

	return item, nil
}
