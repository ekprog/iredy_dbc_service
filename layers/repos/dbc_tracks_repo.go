package repos

import (
	"database/sql"
	"microservice/app/core"
	"microservice/layers/domain"
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
			   ON CONFLICT(user_id, "date") DO UPDATE SET done=$4, updated_at=now();`
	_, err := r.db.Exec(query, track.UserId, track.ChallengeId, track.Date, track.Done)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCTracksRepo) FindByDate(challengeId int32, t time.Time) (done bool, err error) {
	query := `select done from dbc_challenges_tracks 
            		where challenge_id=$1 and date=$2
            		limit 1`

	err = r.db.QueryRow(query, challengeId, t).Scan(&done)
	switch err {
	case nil:
		return done, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}
