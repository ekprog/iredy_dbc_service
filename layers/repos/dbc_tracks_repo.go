package repos

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"microservice/app/core"
	"microservice/layers/domain"
	"strings"
	"time"
)

type DBCTracksRepo struct {
	log core.Logger
	db  *sql.DB
}

func NewDBCTracksRepo(log core.Logger, db *sql.DB) *DBCTracksRepo {
	return &DBCTracksRepo{log: log, db: db}
}

func (r *DBCTracksRepo) InsertOrUpdate(track *domain.DBCTrack) error {
	query := `INSERT INTO dbc_challenges_tracks (user_id, challenge_id, "date", done) 
					VALUES ($1, $2, $3, $4) 
			   ON CONFLICT(challenge_id, "date") DO UPDATE SET done=$4, updated_at=now();`
	_, err := r.db.Exec(query, track.UserId, track.ChallengeId, track.Date.UTC(), track.Done)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCTracksRepo) FindByDate(challengeId int32, t time.Time) (done bool, err error) {
	query := `select done from dbc_challenges_tracks 
            		where challenge_id=$1 and date=$2
            		limit 1`

	err = r.db.QueryRow(query, challengeId, t.UTC()).Scan(&done)
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
