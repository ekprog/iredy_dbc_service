package repos

import (
	"database/sql"
	"fmt"
	"microservice/app/core"
	"microservice/domain"
)

type TasksRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewTasksRepo(log core.Logger, db *sql.DB) *TasksRepo {
	return &TasksRepo{log: log, db: db}
}

func (r *TasksRepo) FetchAll(userId int32, done bool, offset, limit int32) ([]*domain.Task, error) {

	orderBy := "priority desc, created_at desc"
	if done {
		orderBy = "updated_at desc, priority desc"
	}

	query := fmt.Sprintf(`select 
    			id,
    			project_id,
    			name, 
    			"desc", 
    			priority, 
    			done,
    			created_at, 
    			updated_at,
    			deleted_at
			from tasks
			where user_id=$1 and done=$2 and deleted_at is null
			order by %s
			offset $3
			limit $4`, orderBy)

	rows, err := r.db.Query(query, userId, done, offset, limit)
	if err != nil {
		return nil, err
	}

	var result []*domain.Task
	for rows.Next() {
		item := &domain.Task{
			UserId: userId,
		}
		err := rows.Scan(&item.Id,
			&item.ProjectId,
			&item.Name,
			&item.Desc,
			&item.Priority,
			&item.Done,
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

func (r *TasksRepo) FetchById(id int32) (*domain.Task, error) {
	var task = &domain.Task{
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
		&task.ProjectId,
		&task.Name,
		&task.Desc,
		&task.Priority,
		&task.Done,
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

func (r *TasksRepo) FetchByProjectId(projectId int32, done bool, offset, limit int32) ([]*domain.Task, error) {
	orderBy := "priority desc, created_at desc"
	if done {
		orderBy = "created_at desc, priority desc"
	}

	query := `select 
    			id,
    			user_id,
    			name, 
    			"desc", 
    			priority, 
    			done,
    			created_at, 
    			updated_at,
    			deleted_at
			from tasks
			where project_id=$1 and done=$2 and deleted_at is null
			order by $3
			offset $4
			limit $5`

	rows, err := r.db.Query(query, projectId, done, orderBy, offset, limit)
	if err != nil {
		return nil, err
	}

	var result []*domain.Task
	for rows.Next() {
		item := &domain.Task{
			ProjectId: &projectId,
		}
		err := rows.Scan(&item.Id,
			&item.UserId,
			&item.Name,
			&item.Desc,
			&item.Priority,
			&item.Done,
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

func (r *TasksRepo) Insert(item *domain.Task) error {
	query := `INSERT INTO tasks (
                   user_id, 
                   project_id, 
                   name, 
                   "desc", 
                   priority, 
                   done) 
			 VALUES ($1, $2, $3, $4, $5, $6) returning id;`
	err := r.db.QueryRow(query,
		item.UserId,
		item.ProjectId,
		item.Name,
		item.Desc,
		item.Priority,
		item.Done,
	).Scan(&item.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TasksRepo) Update(item *domain.Task) error {
	query := `UPDATE tasks 
				SET name=$2, "desc"=$3, priority=$4, updated_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query,
		item.Id,
		item.Name,
		item.Desc,
		item.Priority)
	if err != nil {
		return err
	}
	return nil
}

func (r *TasksRepo) Remove(id int32) error {
	query := `UPDATE tasks 
				SET deleted_at=now()
				WHERE id=$1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
