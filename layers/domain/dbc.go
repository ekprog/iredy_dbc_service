package domain

import (
	"context"
	"time"
)

//
// MODELS
//

type DBCCategory struct {
	Id     int64
	UserId int64

	Name string

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCChallenge struct {
	Id           int64
	UserId       int64
	CategoryId   *int64
	CategoryName *string
	IsAutoTrack  bool

	Name       string
	Desc       *string
	Image      *string
	LastSeries int64

	LastTracks []*DBCTrack

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCTrack struct {
	Id          int64
	UserId      int64
	ChallengeId int64
	Date        time.Time
	Done        bool

	LastSeries int64
	Score      int64
}

//
// REPOSITORIES
//

type DBCTrackRepository interface {
	InsertOrUpdate(context.Context, *DBCTrack) error
	GetByDate(ctx context.Context, challengeId int64, t time.Time) (*DBCTrack, error)
	FetchForChallengeByDates(challengeId int64, list []time.Time) ([]*DBCTrack, error)
	FetchNotProcessed(challengeId int64, timeSince time.Time) ([]*DBCTrack, error)
	SetProcessed(ctx context.Context, trackIds []int64) error
	Count(ctx context.Context, challengeId int64) (int64, error)
	GetLastForChallengeBefore(ctx context.Context, challengeId int64, date time.Time) (*DBCTrack, error)
	InsertNew(ctx context.Context, tracks []*DBCTrack) error
}

type DBCCategoryRepository interface {
	FetchNotEmptyByUserId(int64) ([]*DBCCategory, error)
	FetchByName(int64, string) (*DBCCategory, error)
	FetchById(int64) (*DBCCategory, error)
	Insert(*DBCCategory) error
	Update(*DBCCategory) error
	Remove(int64, int64) error
}

type DBCChallengesRepository interface {
	FetchUsersAll(userId int64) ([]*DBCChallenge, error)
	FetchAll(limit, offset int64) ([]*DBCChallenge, error)
	FetchById(context.Context, int64) (*DBCChallenge, error)
	FetchByName(int64, string) (*DBCChallenge, error)
	Insert(*DBCChallenge) error
	Update(*DBCChallenge) error
	Remove(int64) error
}

//
// USE CASES
//

type DBCCategoryUseCase interface {
	Get(userId int64) (CategoryListResponse, error)
	Update(*DBCCategory) (StatusResponse, error)
	Remove(userId, taskId int64) (StatusResponse, error)
}

type DBCChallengesUseCase interface {
	All(userId int64) (ChallengesListResponse, error)
	Create(form *CreateDBCChallengeForm) (CreateChallengeResponse, error)
	Update(ctx context.Context, task *DBCChallenge) (StatusResponse, error)
	Remove(userId, taskId int64) (StatusResponse, error)

	TrackDay(ctx context.Context, form *DBCTrack) (UserGamifyResponse, error)
}

// IO FORMS (FORMS)

type CreateDBCChallengeForm struct {
	UserId       int64
	Name         string
	Desc         *string
	CategoryName *string
	IsAutoTrack  bool
}

// IO FORMS (RESPONSES)

type CreateChallengeResponse struct {
	StatusCode string
	Id         int64
	CategoryId *int64
}

type ChallengesListResponse struct {
	StatusCode string
	Challenges []*DBCChallenge
}

type CategoryListResponse struct {
	StatusCode string
	Categories []*DBCCategory
}
