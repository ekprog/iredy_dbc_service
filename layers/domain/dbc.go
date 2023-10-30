package domain

import (
	"context"
	"time"
)

//
// MODELS
//

type DBCCategory struct {
	Id        int64
	UserId    int64
	Name      string
	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCChallengeInfo struct {
	Id      int64
	OwnerId int64

	CategoryId *int64
	Category   *DBCCategory

	IsAutoTrack    bool
	VisibilityType string

	Name  string
	Desc  *string
	Image *string

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCUserChallenge struct {
	Id int64

	ChallengeInfoId int64
	ChallengeInfo   *DBCChallengeInfo

	UserId     int64
	LastSeries int64
	LastTracks []*DBCTrack

	UpdatedAt time.Time
	CreatedAt time.Time
	DeletedAt *time.Time
}

type DBCTrack struct {
	Id     int64
	UserId int64

	ChallengeUserId int64
	ChallengeId     int64

	Date time.Time
	Done bool

	LastSeries int64
	Score      int64
	ScoreDaily int64
}

// REPOSITORIES
type DBCCategoryRepository interface {
	FetchNotEmptyByUserId(int64) ([]*DBCCategory, error)
	FetchByName(int64, string) (*DBCCategory, error)
	FetchById(int64) (*DBCCategory, error)
	Insert(*DBCCategory) error
	Update(*DBCCategory) error
	Remove(int64, int64) error
}

type DBChallengeInfoRepository interface {
	// No scope
	FetchById(int64) (*DBCChallengeInfo, error)
	Insert(item *DBCChallengeInfo) error

	// Public scope
	PublicFetchLike(search string, categoryId *int64, limit, offset int64) ([]*DBCChallengeInfo, error)
}

type DBCUserChallengeRepository interface {
	// No scope
	FetchAll(limit, offset int64) ([]*DBCUserChallenge, error)
	FetchById(context.Context, int64) (*DBCUserChallenge, error)
	Insert(*DBCUserChallenge) error
	Update(*DBCUserChallenge) error
	Remove(int64) error

	// User scope
	UserFetchAll(userId int64) ([]*DBCUserChallenge, error)
	UserFetchByName(int64, string) (*DBCUserChallenge, error)
	UserExistsByChallengeId(int64, int64) (bool, error)
}

type DBCTrackRepository interface {
	// No scope
	SetProcessed(ctx context.Context, trackIds []int64) error
	InsertOrUpdateBulk(context.Context, []*DBCTrack) error

	// Challenge scope
	ChallengeFetchByDates(challengeId int64, list []time.Time) ([]*DBCTrack, error)
	ChallengeFetchLastBefore(ctx context.Context, challengeId int64, date time.Time) (*DBCTrack, error)
	ChallengeFetchLast(ctx context.Context, challengeId int64) (*DBCTrack, error)
	ChallengeFetchAfter(ctx context.Context, challengeId int64, date time.Time) ([]*DBCTrack, error)
	ChallengeFetchBetween(ctx context.Context, challengeId int64, from, to time.Time) ([]*DBCTrack, error)

	// Challenge Not processed scope
	NotProcessedChallengeFetchAllBefore(ctx context.Context, challengeId int64, date time.Time) ([]*DBCTrack, error)
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
	// Public Scope
	PublicSearch(search string, categoryId *int64, limit, offset int64) (ChallengesListResponse, error)

	// User scope
	UserAll(userId int64) (UserChallengesListResponse, error)
	UserCreate(form *CreateDBCChallengeForm) (CreateChallengeResponse, error)

	//
	Info(userId int64, id int64) (ChallengeInfoResponse, error)
	Update(ctx context.Context, task *DBCUserChallenge) (StatusResponse, error)
	Remove(userId, taskId int64) (StatusResponse, error)

	TrackDay(ctx context.Context, form *DBCTrack) (UserGamifyResponse, error)
	GetMonthTracks(ctx context.Context, date time.Time, challengeId, userId int64) (*ChallengeMonthTracksResponse, error)
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

type UserChallengesListResponse struct {
	StatusCode     string
	UserChallenges []*DBCUserChallenge
}

type ChallengeInfoResponse struct {
	StatusCode string
	Challenge  *DBCChallengeInfo
	IsMember   bool
}

type ChallengesListResponse struct {
	StatusCode string
	Challenges []*DBCChallengeInfo
}

type CategoryListResponse struct {
	StatusCode string
	Categories []*DBCCategory
}

type ChallengeMonthTracksResponse struct {
	StatusCode string
	Tracks     []*DBCTrack
}
