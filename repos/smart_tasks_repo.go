package repos

import (
	"database/sql"
	"fmt"
	"microservice/app/core"
	"microservice/domain"
	"time"
)

type SmartTasksRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewSmartTasksRepo(log core.Logger, db *sql.DB) *SmartTasksRepo {
	return &SmartTasksRepo{log: log, db: db}
}

func (r *SmartTasksRepo) CountAllActive() (count int32, err error) {
	query := `select count(id)
				from smart_tasks
				where deleted_at is null;`

	err = r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (r *SmartTasksRepo) FetchAllActive(offset, limit int32) ([]*domain.SmartTask, error) {

	query := `select st.id,
       			   st.user_id,
				   st.project_id,
				   st.name,
				   st.desc,
				   st.priority,
				   st.last_generated_at at time zone 'Europe/Moscow',
				   st.created_at at time zone 'Europe/Moscow',
				   st.updated_at at time zone 'Europe/Moscow',
				   st.deleted_at at time zone 'Europe/Moscow',
				   stg.id,
				   stg.period,
				   stg.datetime at time zone 'Europe/Moscow',
				   stg.created_at at time zone 'Europe/Moscow',
				   stg.updated_at at time zone 'Europe/Moscow',
				   stg.deleted_at at time zone 'Europe/Moscow'
			from smart_tasks st
					 left join smart_tasks_gen stg
							   on st.id = stg.smart_task_id
			WHERE st.deleted_at is null and
			      stg.deleted_at is null and
			      st.id IN (SELECT id
							FROM smart_tasks
							LIMIT $1 OFFSET $2);`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	result := map[int32]*domain.SmartTask{}
	for rows.Next() {
		itemP := &domain.SmartTask{}

		// For generation item (Can be null)
		var id *int32
		var period *string
		var datetime, createdAt, updatedAt, deletedAt *time.Time

		// Scan
		err := rows.Scan(&itemP.Id,
			&itemP.UserId,
			&itemP.ProjectId,
			&itemP.Name,
			&itemP.Desc,
			&itemP.Priority,
			&itemP.LastGeneratedAt,
			&itemP.CreatedAt,
			&itemP.UpdatedAt,
			&itemP.DeletedAt,
			&id,
			&period,
			&datetime,
			&createdAt,
			&updatedAt,
			&deletedAt)
		if err != nil {
			return nil, err
		}

		// Check if doesnt exists in Result
		if _, ok := result[itemP.Id]; !ok {
			result[itemP.Id] = itemP
		}

		// Generation Item
		if id != nil {
			itemC := &domain.GenerationItem{
				Id:        *id,
				Period:    domain.GenerationPeriod(*period),
				Datetime:  *datetime,
				CreatedAt: *createdAt,
				UpdatedAt: *updatedAt,
				DeletedAt: deletedAt,
			}
			itemC.Period = domain.GenerationPeriod(*period)
			result[itemP.Id].GenerationItems = append(result[itemP.Id].GenerationItems, itemC)
		}

	}

	// To array
	values := make([]*domain.SmartTask, 0, len(result))
	for _, v := range result {
		values = append(values, v)
	}

	return values, nil
}

func (r *SmartTasksRepo) FetchAllByUserId(userId int32, trashed bool, offset, limit int32) ([]*domain.SmartTask, error) {

	deletedAtQuery := "is null"
	if trashed {
		deletedAtQuery = "is not null"
	}

	query := fmt.Sprintf(`select st.id,
				   st.project_id,
				   st.name,
				   st.desc,
				   st.priority,
				   st.last_generated_at,
				   st.created_at,
				   st.updated_at,
				   st.deleted_at,
				   stg.id,
				   stg.period,
				   stg.datetime,
				   stg.created_at,
				   stg.updated_at,
				   stg.deleted_at
			from smart_tasks st
					 left join smart_tasks_gen stg
							   on st.id = stg.smart_task_id
			
			WHERE st.user_id=$1 and
			      st.deleted_at %s and
			      stg.deleted_at is null and
			      st.id IN (SELECT id
							FROM smart_tasks
							LIMIT $2 OFFSET $3);`, deletedAtQuery)

	rows, err := r.db.Query(query, userId, limit, offset)
	if err != nil {
		return nil, err
	}

	result := map[int32]*domain.SmartTask{}
	for rows.Next() {
		itemP := &domain.SmartTask{
			UserId: userId,
		}

		// For generation item (Can be null)
		var id *int32
		var period *string
		var datetime, createdAt, updatedAt, deletedAt *time.Time

		// Scan
		err := rows.Scan(&itemP.Id,
			&itemP.ProjectId,
			&itemP.Name,
			&itemP.Desc,
			&itemP.Priority,
			&itemP.LastGeneratedAt,
			&itemP.CreatedAt,
			&itemP.UpdatedAt,
			&itemP.DeletedAt,
			&id,
			&period,
			&datetime,
			&createdAt,
			&updatedAt,
			&deletedAt)
		if err != nil {
			return nil, err
		}

		// Check if doesnt exists in Result
		if _, ok := result[itemP.Id]; !ok {
			result[itemP.Id] = itemP
		}

		// Generation Item
		if id != nil {
			itemC := &domain.GenerationItem{
				Id:        *id,
				Period:    domain.GenerationPeriod(*period),
				Datetime:  *datetime,
				CreatedAt: *createdAt,
				UpdatedAt: *updatedAt,
				DeletedAt: deletedAt,
			}
			itemC.Period = domain.GenerationPeriod(*period)
			result[itemP.Id].GenerationItems = append(result[itemP.Id].GenerationItems, itemC)
		}

	}

	// To array
	values := make([]*domain.SmartTask, 0, len(result))
	for _, v := range result {
		values = append(values, v)
	}

	return values, nil
}

func (r *SmartTasksRepo) UpdateLastGeneratedAt(id int32, time time.Time) error {
	query := `UPDATE smart_tasks 
				SET last_generated_at=$2
				WHERE id=$1`
	_, err := r.db.Exec(query, id, time)
	if err != nil {
		return err
	}
	return nil
}
