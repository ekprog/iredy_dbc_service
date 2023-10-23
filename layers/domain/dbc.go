package domain

import (
	"context"
	"time"
)

//
// MODELS
//

type DBCCategory struct {
	Id     int32
	UserId int32

	Name string

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCChallenge struct {
	Id           int32
	UserId       int32
	CategoryId   *int32
	CategoryName *string
	IsAutoTrack  bool

	Name       string
	Desc       *string
	Image      *string
	LastSeries int32

	LastTracks []*DBCTrack

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCTrack struct {
	Id          int64
	UserId      int32
	ChallengeId int32
	Date        time.Time
	Done        bool

	LastSeries int32
	Score      int32
}

//
// REPOSITORIES
//

type DBCTrackRepository interface {
	InsertOrUpdate(context.Context, *DBCTrack) error
	CheckDoneByDate(ctx context.Context, challengeId int32, t time.Time) (bool, error)
	FetchForChallengeByDates(challengeId int32, list []time.Time) ([]*DBCTrack, error)
	FetchNotProcessed(challengeId int32, timeSince time.Time) ([]*DBCTrack, error)
	SetProcessed(ctx context.Context, trackIds []int64) error
}

type DBCCategoryRepository interface {
	FetchNotEmptyByUserId(int32) ([]*DBCCategory, error)
	FetchByName(int32, string) (*DBCCategory, error)
	FetchById(int32) (*DBCCategory, error)
	Insert(*DBCCategory) error
	Update(*DBCCategory) error
	Remove(int32, int32) error
}

type DBCChallengesRepository interface {
	FetchUsersAll(userId int32) ([]*DBCChallenge, error)
	FetchAll(limit, offset int64) ([]*DBCChallenge, error)
	FetchById(int32) (*DBCChallenge, error)
	FetchByName(int32, string) (*DBCChallenge, error)
	Insert(*DBCChallenge) error
	Update(*DBCChallenge) error
	Remove(int32) error
}

//
// USE CASES
//

type DBCCategoryUseCase interface {
	Get(userId int32) (CategoryListResponse, error)
	Update(*DBCCategory) (StatusResponse, error)
	Remove(userId, taskId int32) (StatusResponse, error)
}

type DBCChallengesUseCase interface {
	All(userId int32) (ChallengesListResponse, error)
	Create(form *CreateDBCChallengeForm) (CreateChallengeResponse, error)
	Update(task *DBCChallenge) (StatusResponse, error)
	Remove(userId, taskId int32) (StatusResponse, error)

	TrackDay(ctx context.Context, form *DBCTrack) (UserGamifyResponse, error)
}

// IO FORMS (FORMS)

type CreateDBCChallengeForm struct {
	UserId       int32
	Name         string
	Desc         *string
	CategoryName *string
	IsAutoTrack  bool
}

// IO FORMS (RESPONSES)

type CreateChallengeResponse struct {
	StatusCode string
	Id         int32
	CategoryId *int32
}

type ChallengesListResponse struct {
	StatusCode string
	Challenges []*DBCChallenge
}

type CategoryListResponse struct {
	StatusCode string
	Categories []*DBCCategory
}
