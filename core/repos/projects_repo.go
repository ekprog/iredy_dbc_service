package repos

import (
	"database/sql"
	"fmt"
	"github.com/samber/lo"
	"microservice/app/core"
	"microservice/core/domain"
	"strings"
)

type ProjectsRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewProjectsRepo(log core.Logger, db *sql.DB) *ProjectsRepo {
	return &ProjectsRepo{log: log, db: db}
}

func (r *ProjectsRepo) FetchByUserId(userId int32) ([]*domain.Project, error) {
	query := `select 
    			id, 
    			name,
    			"desc", 
    			color,
    			"order",
    			parent_id,
    			created_at,
    			updated_at,
    			deleted_at
			from projects
			where user_id=$1 and deleted_at is null`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var result []*domain.Project
	for rows.Next() {
		item := &domain.Project{
			UserId: userId,
		}
		err := rows.Scan(&item.Id,
			&item.Name,
			&item.Desc,
			&item.Color,
			&item.Order,
			&item.ParentId,
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

func (r *ProjectsRepo) FetchByUserIdTrashed(userId int32) ([]*domain.Project, error) {
	query := `select 
    			id, 
    			name,
    			"desc", 
    			color,
    			"order",
    			parent_id,
    			created_at,
    			updated_at,
    			deleted_at
			from projects
			where user_id=$1 and deleted_at is not null`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var result []*domain.Project
	for rows.Next() {
		item := &domain.Project{
			UserId: userId,
		}
		err := rows.Scan(
			&item.Id,
			&item.Name,
			&item.Desc,
			&item.Color,
			&item.Order,
			&item.ParentId,
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

func (r *ProjectsRepo) FetchById(id int32) (*domain.Project, error) {
	var item = &domain.Project{
		Id: id,
	}
	query := `select 
    			user_id,
    			name,
    			"desc", 
    			color,
    			"order",
    			parent_id,
    			created_at,
    			updated_at,
    			deleted_at
			from projects
			where id=$1
			limit 1`

	err := r.db.QueryRow(query, id).Scan(
		&item.UserId,
		&item.Name,
		&item.Desc,
		&item.Color,
		&item.Order,
		&item.ParentId,
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

func (r *ProjectsRepo) Insert(item *domain.Project) error {
	query := `INSERT INTO projects (user_id, name, "desc", color, "order", parent_id) VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	err := r.db.QueryRow(query, item.UserId, item.Name, item.Desc, item.Color, item.Order, item.ParentId).Scan(&item.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectsRepo) Update(item *domain.Project) error {
	query := `UPDATE projects 
				SET name=$2, "desc"=$3, color=$4, updated_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, item.Id, item.Name, item.Desc, item.Color)
	if err != nil {
		return err
	}
	return nil
}

// Обновляет порядок проектов для пользователя
func (r *ProjectsRepo) UpdateOrderForUser(userId int32, itemsIds []int32) error {

	idOrder := lo.Map(itemsIds, func(x int32, index int) string {
		return fmt.Sprintf("(%d, %d)", x, index)
	})
	idOrderSql := strings.Join(idOrder, ",")

	query := fmt.Sprintf(`UPDATE projects 
				SET "order"=t."order"
				FROM ( VALUES %s) as t (id, "order")
				WHERE projects.id = t.id AND projects.user_id=$1`, idOrderSql)
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectsRepo) DragItemsForUser(userId int32, drags []*domain.ProjectDrag) error {

	idOrder := lo.Map(drags, func(x *domain.ProjectDrag, index int) string {
		parentId := int32(-1)
		if x.ParentId != nil {
			parentId = *x.ParentId
		}
		return fmt.Sprintf("(%d, %d, %d)", x.ProjectId, parentId, x.Order)
	})
	idOrderSql := strings.Join(idOrder, ",")

	query := fmt.Sprintf(`UPDATE projects 
				SET "order"=t."order", parent_id=case when t.parent_id=-1 then null else t.parent_id end, updated_at=now()
				FROM ( VALUES %s) as t (id, parent_id, "order")
				WHERE projects.id = t.id AND projects.user_id=$1`, idOrderSql)
	_, err := r.db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ProjectsRepo) Remove(id int32) error {
	query := `UPDATE projects 
				SET deleted_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
