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
	"microservice/tools"
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

func (r *DBCTracksRepo) ChallengeFetchLastBefore(ctx context.Context, challengeId int64, date time.Time) (track *domain.DBCTrack, err error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				challenge_id,
    				"date",
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenge_tracks 
            		where challenge_user_id=$1 and "date" < $2
            		order by "date" desc
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId, date).Scan(
		&track.Id,
		&track.UserId,
		&track.ChallengeId,
		&track.Date,
		&track.Done,
		&track.LastSeries,
		&track.Score,
		&track.ScoreDaily)
	switch err {
	case nil:
		return track, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "ChallengeFetchLastBefore")
	}
}

func (r *DBCTracksRepo) GetLastNotProcessedForChallengeBefore(ctx context.Context, challengeId int64, date time.Time) (track *domain.DBCTrack, err error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				challenge_id,
    				date,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenge_tracks 
            		where challenge_id=$1 and "date" < $2 and processed = false
            		order by "date" desc
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId, date).Scan(
		&track.Id,
		&track.UserId,
		&track.ChallengeId,
		&track.Date,
		&track.Done,
		&track.LastSeries,
		&track.Score,
		&track.ScoreDaily)
	switch err {
	case nil:
		return track, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "ChallengeFetchLastBefore")
	}
}

func (r *DBCTracksRepo) ChallengeFetchLast(ctx context.Context, challengeId int64) (track *domain.DBCTrack, err error) {
	query := `select 
    				id,
    				user_id,
    				challenge_id,
    				date,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenge_tracks 
            		where challenge_user_id=$1
            		order by "date" desc
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId).Scan(
		&track.Id,
		&track.UserId,
		&track.ChallengeId,
		&track.Date,
		&track.Done,
		&track.LastSeries,
		&track.Score,
		&track.ScoreDaily)
	switch err {
	case nil:
		return track, nil
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, errors.Wrap(err, "ChallengeFetchLast")
	}
}

func (r *DBCTracksRepo) ChallengeFetchByDates(challengeUserId int64, list []time.Time) ([]*domain.DBCTrack, error) {

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
      						from dbc_challenge_tracks t
               					right join (select date
                           			from (values %s) s(date)) s
                          			on s.date = t.date and challenge_id = $1
							order by t.date asc) st`, strings.Join(dateStrings, ","))

	rows, err := r.db.Query(query, challengeUserId)
	if err != nil {
		return nil, errors.Wrap(err, "ChallengeFetchByDates")
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

func (r *DBCTracksRepo) ChallengeFetchAfter(ctx context.Context, challengeUserId int64, date time.Time) ([]*domain.DBCTrack, error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				challenge_id,
    				"date",
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenge_tracks 
            		where challenge_user_id=$1 and "date" > $2
            		order by "date"`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).QueryContext(ctx, query, challengeUserId, date)
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCTrack
	for rows.Next() {
		item := &domain.DBCTrack{
			ChallengeUserId: challengeUserId,
		}
		err := rows.Scan(
			&item.Id,
			&item.UserId,
			&item.ChallengeId,
			&item.Date,
			&item.Done,
			&item.LastSeries,
			&item.Score,
			&item.ScoreDaily)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *DBCTracksRepo) ChallengeFetchBetween(ctx context.Context, challengeUserId int64, from, to time.Time) ([]*domain.DBCTrack, error) {
	from = tools.RoundDateTimeToDay(from.UTC())
	to = tools.RoundDateTimeToDay(to.UTC())

	query := `select 
    				id,
    				user_id,
    				challenge_id,
    				"date",
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenge_tracks 
            		where challenge_user_id=$1 and 
            		      "date" >= $2 and "date" <= $3
            		order by "date"`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).QueryContext(ctx, query, challengeUserId, from, to)
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCTrack
	for rows.Next() {
		item := &domain.DBCTrack{
			ChallengeUserId: challengeUserId,
		}
		err := rows.Scan(
			&item.Id,
			&item.UserId,
			&item.ChallengeId,
			&item.Date,
			&item.Done,
			&item.LastSeries,
			&item.Score,
			&item.ScoreDaily)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (r *DBCTracksRepo) NotProcessedChallengeFetchAllBefore(ctx context.Context, challengeUserId int64, date time.Time) ([]*domain.DBCTrack, error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				challenge_id,
    				"date",
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenge_tracks 
            		where challenge_user_id=$1 and "date" < $2 and processed = false
            		order by "date"`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).QueryContext(ctx, query, challengeUserId, date)
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCTrack
	for rows.Next() {
		item := &domain.DBCTrack{
			ChallengeUserId: challengeUserId,
		}
		err := rows.Scan(
			&item.Id,
			&item.UserId,
			&item.ChallengeId,
			&item.Date,
			&item.Done,
			&item.LastSeries,
			&item.Score,
			&item.ScoreDaily)
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

	query := `UPDATE dbc_challenge_tracks 
				SET processed=true, updated_at=now()
				where id in (%s)`

	query = fmt.Sprintf(query, inParams)

	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCTracksRepo) InsertOrUpdateBulk(ctx context.Context, tracks []*domain.DBCTrack) error {

	if len(tracks) == 0 {
		return nil
	}

	var userIdList, challengeIdList, challengeUserIdList, lastSeriesList, scoreList, scoreDailyList, doneList, dateList []string

	lo.ForEach(tracks, func(track *domain.DBCTrack, index int) {
		userIdList = append(userIdList, fmt.Sprintf("%d", track.UserId))
		challengeIdList = append(challengeIdList, fmt.Sprintf("%d", track.ChallengeId))
		challengeUserIdList = append(challengeUserIdList, fmt.Sprintf("%d", track.ChallengeUserId))
		lastSeriesList = append(lastSeriesList, fmt.Sprintf("%d", track.LastSeries))
		scoreList = append(scoreList, fmt.Sprintf("%d", track.Score))
		scoreDailyList = append(scoreDailyList, fmt.Sprintf("%d", track.ScoreDaily))
		doneList = append(doneList, fmt.Sprintf("%v", track.Done))
		dateList = append(dateList, fmt.Sprintf("'%s'::date", track.Date.Format("2006-01-02")))
	})

	query := fmt.Sprintf(`
		insert into dbc_challenge_tracks (
			user_id,
			challenge_id,
			challenge_user_id,
			"date",
			done,
			last_series,
			score,
			score_daily
		)
		select unnest(array[%s]),
			   unnest(array[%s]),
			   unnest(array[%s]),
			   unnest(array[%s]),
			   unnest(array[%s]),
			   unnest(array[%s]),
			   unnest(array[%s]),
				unnest(array[%s])
		on conflict (challenge_id, "date") do
			update set
					   "date" = excluded.date,
					   score = excluded.score,
					   score_daily = excluded.score_daily,
					   last_series = excluded.last_series,
					   done = excluded.done, 
					   updated_at=now()`,
		strings.Join(userIdList, ","),
		strings.Join(challengeIdList, ","),
		strings.Join(challengeUserIdList, ","),
		strings.Join(dateList, ","),
		strings.Join(doneList, ","),
		strings.Join(lastSeriesList, ","),
		strings.Join(scoreList, ","),
		strings.Join(scoreDailyList, ","),
	)

	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil

}
