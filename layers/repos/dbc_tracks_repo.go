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

func (r *DBCTracksRepo) InsertOrUpdate(ctx context.Context, track *domain.DBCTrack) error {
	query := `INSERT INTO dbc_challenges_tracks (user_id, challenge_id, "date", done, last_series, score, score_daily) 
					VALUES ($1, $2, $3, $4, $5, $6, $7) 
			   ON CONFLICT(challenge_id, "date") 
			       DO UPDATE SET done=$4, last_series=$5, score=$6, score_daily=$7, updated_at=now();`
	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query,
		track.UserId,
		track.ChallengeId,
		track.Date.UTC(),
		track.Done,
		track.LastSeries,
		track.Score,
		track.ScoreDaily)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCTracksRepo) Count(ctx context.Context, challengeId int64) (c int64, err error) {
	query := `select  count(id) 
				from dbc_challenges_tracks 
            	where challenge_id=$1`

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId).Scan(&c)
	if err != nil {
		return -1, errors.Wrap(err, "Count")
	}
	return c, nil
}

func (r *DBCTracksRepo) GetLastForChallengeBefore(ctx context.Context, challengeId int64, date time.Time) (track *domain.DBCTrack, err error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				date,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenges_tracks 
            		where challenge_id=$1 and "date" < $2
            		order by "date" desc
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId, date).Scan(
		&track.Id,
		&track.UserId,
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
		return nil, errors.Wrap(err, "GetLastForChallengeBefore")
	}
}

func (r *DBCTracksRepo) GetLastNotProcessedForChallengeBefore(ctx context.Context, challengeId int64, date time.Time) (track *domain.DBCTrack, err error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				date,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenges_tracks 
            		where challenge_id=$1 and "date" < $2 and processed = false
            		order by "date" desc
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId, date).Scan(
		&track.Id,
		&track.UserId,
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
		return nil, errors.Wrap(err, "GetLastForChallengeBefore")
	}
}

func (r *DBCTracksRepo) GetLastForChallenge(ctx context.Context, challengeId int64) (track *domain.DBCTrack, err error) {
	query := `select 
    				id,
    				user_id,
    				date,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenges_tracks 
            		where challenge_id=$1
            		order by "date" desc
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId).Scan(
		&track.Id,
		&track.UserId,
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
		return nil, errors.Wrap(err, "GetLastForChallenge")
	}
}

func (r *DBCTracksRepo) GetByDate(ctx context.Context, challengeId int64, date time.Time) (track *domain.DBCTrack, err error) {
	date = tools.RoundDateTimeToDay(date)

	query := `select 
    				id,
    				user_id,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenges_tracks 
            		where challenge_id=$1 and date=$2
            		limit 1`

	track = &domain.DBCTrack{
		ChallengeId: challengeId,
		Date:        date,
	}

	err = r.getter.DefaultTrOrDB(ctx, r.db).QueryRowContext(ctx, query, challengeId, date.UTC()).Scan(
		&track.Id,
		&track.UserId,
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
		return nil, errors.Wrap(err, "GetByDate")
	}
}

func (r *DBCTracksRepo) FetchForChallengeByDates(challengeId int64, list []time.Time) ([]*domain.DBCTrack, error) {

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

func (r *DBCTracksRepo) GetAllForChallengeAfter(ctx context.Context, challengeId int64, date time.Time) ([]*domain.DBCTrack, error) {
	date = tools.RoundDateTimeToDay(date.UTC())

	query := `select 
    				id,
    				user_id,
    				date,
    				done, 
       				last_series, 
       				score,
       				score_daily from dbc_challenges_tracks 
            		where challenge_id=$1 and "date" > $2
            		order by "date"`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.db).QueryContext(ctx, query, challengeId, date)
	if err != nil {
		return nil, err
	}

	var result []*domain.DBCTrack
	for rows.Next() {
		item := &domain.DBCTrack{}
		err := rows.Scan(&item.Id, &item.UserId, &item.Date, &item.Done, &item.LastSeries, &item.Score, &item.ScoreDaily)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

// Возвращает треки, подлежащие учету в dailyScore
// timeSince - время ДО которого искать (зависит от типа генерации треков в челлендже)
func (r *DBCTracksRepo) FetchNotProcessed(challengeId int64, timeSince time.Time) ([]*domain.DBCTrack, error) {
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

func (r *DBCTracksRepo) InsertNew(ctx context.Context, tracks []*domain.DBCTrack) error {

	if len(tracks) == 0 {
		return nil
	}

	valuesArr := lo.Map(tracks, func(track *domain.DBCTrack, index int) string {
		return fmt.Sprintf(`(%d, %d, '%s'::date, %v, %d, %d, %d)`,
			track.UserId,
			track.ChallengeId,
			track.Date.Format("2006-01-02"),
			track.Done,
			track.LastSeries,
			track.Score,
			track.ScoreDaily,
		)
	})
	values := strings.Join(valuesArr, ",")

	query := fmt.Sprintf(`INSERT INTO dbc_challenges_tracks (
                                   user_id, 
                                   challenge_id, 
                                   date, 
                                   done, 
                                   last_series, 
                                   score,
                                   score_daily) VALUES %s`, values)

	query = fmt.Sprintf(query)

	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, query)
	if err != nil {
		return err
	}
	return nil
}

func (r *DBCTracksRepo) UpdateSome(ctx context.Context, tracks []*domain.DBCTrack) error {

	if len(tracks) == 0 {
		return nil
	}

	valuesArr := lo.Map(tracks, func(track *domain.DBCTrack, index int) string {
		return fmt.Sprintf(`(%d, %d, %d, %d)`,
			track.Id,
			track.LastSeries,
			track.Score,
			track.ScoreDaily,
		)
	})
	values := strings.Join(valuesArr, ",")

	query := `update dbc_challenges_tracks as t set
				score = c.score, score_daily=c.score_daily, last_series = c.last_series, updated_at=now() 
		  	 from (values %s) as c(id, last_series, score, score_daily)
			 where c.id = t.id`
	query = fmt.Sprintf(query, values)

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

	var userIdList, challengeIdList, lastSeriesList, scoreList, scoreDailyList, doneList, dateList []string

	lo.ForEach(tracks, func(track *domain.DBCTrack, index int) {
		userIdList = append(userIdList, fmt.Sprintf("%d", track.UserId))
		challengeIdList = append(challengeIdList, fmt.Sprintf("%d", track.ChallengeId))
		lastSeriesList = append(lastSeriesList, fmt.Sprintf("%d", track.LastSeries))
		scoreList = append(scoreList, fmt.Sprintf("%d", track.Score))
		scoreDailyList = append(scoreDailyList, fmt.Sprintf("%d", track.ScoreDaily))
		doneList = append(doneList, fmt.Sprintf("%v", track.Done))
		dateList = append(dateList, fmt.Sprintf("'%s'::date", track.Date.Format("2006-01-02")))
	})

	query := fmt.Sprintf(`
		insert into dbc_challenges_tracks (
			user_id,
			challenge_id,
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
				unnest(array[%s])
		on conflict (challenge_id, "date") do
			update set
					   "date" = excluded.date,
					   score = excluded.score,
					   score_daily = excluded.score_daily,
					   done = excluded.done, 
					   updated_at=now()`,
		strings.Join(userIdList, ","),
		strings.Join(challengeIdList, ","),
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
