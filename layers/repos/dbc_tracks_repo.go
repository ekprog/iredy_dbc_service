package repos

import (
	"context"
	"database/sql"
	"fmt"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"microservice/app/core"
	"microservice/layers/domain"
	"strconv"
	"strings"
	"time"
)

type DBCTracksRepo struct {
	log    core.Logger
	db     *sql.DB
	getter *trmsql.CtxGetter
}

func NewDBCTracksRepo(log core.Logger, db *sql.DB, getter *trmsql.CtxGetter) *DBCTracksRepo {
	return &DBCTracksRepo{
		log:    log,
		db:     db,
		getter: getter,
	}
}

func (r *DBCTracksRepo) InsertOrUpdate(ctx context.Context, track *domain.DBCTrack) error {
	query := `INSERT INTO dbc_challenges_tracks (user_id, challenge_id, "date", done, last_series, score) 
					VALUES ($1, $2, $3, $4, $5, $6) 
			   ON CONFLICT(challenge_id, "date") DO UPDATE SET done=$4, last_series=$5, score=$6, updated_at=now();`
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query,
		track.UserId,
		track.ChallengeId,
		track.Date.UTC(),
		track.Done,
		track.LastSeries,
		track.Score)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCTracksRepo) CheckDoneByDate(ctx context.Context, challengeId int32, t time.Time) (done bool, err error) {
	query := `select done from dbc_challenges_tracks 
            		where challenge_id=$1 and date=$2
            		limit 1`

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId, t.UTC()).Scan(&done)
	switch err {
	case nil:
		return done, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}

func (r *DBCTracksRepo) FetchForChallengeByDates(challengeId int32, list []time.Time) ([]*domain.DBCTrack, error) {

	dateStrings := lo.Map(list, func(item time.Time, index int) string {
		return fmt.Sprintf("('%s'::date)", item.UTC().Format("2006-01-02"))
	})

	query := fmt.Sprintf(`
		select  st.date,
				case
				   when st.done is null then false
				   else st.done
				end as done
					from (select s.date as date, t.done as done
      						from dbc_challenges_tracks t
               					right join (select date
                           			from (values %s) s(date)) s
                          			on s.date = t.date and challenge_id = $1
							order by t.date asc) st`, strings.Join(dateStrings, ","))

	rows, err := r.db.Query(query, challengeId)
	if err != nil {
		return nil, errors.Wrap(err, "FetchForChallengeByDates")
	}

	var result []*domain.DBCTrack
	for rows.Next() {
		item := &domain.DBCTrack{}
		err := rows.Scan(&item.Date, &item.Done)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

// Возвращает треки, подлежащие учету в dailyScore
// timeSince - время ДО которого искать (зависит от типа генерации треков в челлендже)
func (r *DBCTracksRepo) FetchNotProcessed(challengeId int32, timeSince time.Time) ([]*domain.DBCTrack, error) {
	query := `select t.id,
       				t.date,
       				t.done
		  		from dbc_challenges_tracks t
		  			where t.processed = false and 
		  			      t.challenge_id = $1 and
		  			      t.date < $2`

	rows, err := r.db.Query(query, challengeId, timeSince)
	if err != nil {
		return nil, errors.Wrap(err, "FetchNotProcessed")
	}

	var result []*domain.DBCTrack
	for rows.Next() {
		item := &domain.DBCTrack{}
		err := rows.Scan(&item.Id, &item.Date, &item.Done)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *DBCTracksRepo) SetProcessed(ctx context.Context, trackIds []int64) error {

	if len(trackIds) == 0 {
		return nil
	}

	ids := lo.Map(trackIds, func(id int64, index int) string {
		return strconv.FormatInt(id, 10)
	})
	inParams := strings.Join(ids, ",")

	query := `UPDATE dbc_challenges_tracks 
				SET processed=true, updated_at=now()
				where id in (%s)`

	query = fmt.Sprintf(query, inParams)

	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query)
	if err != nil {
		return errors.Wrap(err, "SetProcessed")
	}
	return nil
}
