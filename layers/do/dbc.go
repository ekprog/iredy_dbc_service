package do

import (
	"errors"
	"gorm.io/gorm"
	"microservice/layers/domain"
	"time"
)

type DBCChallengeCategory struct {
	gorm.Model

	UserID int64
	User   *User

	Name string
}

func (m *DBCChallengeCategory) DTO() *domain.DBCCategory {
	if m == nil {
		return nil
	}

	obj := &domain.DBCCategory{
		Id:        int64(m.ID),
		UserId:    m.UserID,
		Name:      m.Name,
		UpdatedAt: m.UpdatedAt,
		CreatedAt: m.CreatedAt,
		DeletedAt: nil,
	}

	if m.DeletedAt.Valid {
		obj.DeletedAt = new(time.Time)
		*obj.DeletedAt = m.DeletedAt.Time
	}

	return obj
}

type DBCChallenge struct {
	gorm.Model

	OwnerID int64
	Owner   *User

	CategoryID int64
	Category   *DBCChallengeCategory

	Name           string
	Image          *string
	Desc           *string
	IsAutoTrack    bool
	VisibilityType string
}

func NewDBCChallenge(from *domain.DBCChallengeInfo) (*DBCChallenge, error) {
	if from == nil {
		return nil, errors.New("from is nil")
	}
	doItem := &DBCChallenge{
		OwnerID:        from.OwnerId,
		Name:           from.Name,
		Image:          from.Image,
		Desc:           from.Desc,
		IsAutoTrack:    from.IsAutoTrack,
		VisibilityType: from.VisibilityType,
	}
	if from.Category != nil {
		doItem.CategoryID = from.Category.Id
	}

	return doItem, nil
}

func (m *DBCChallenge) DTO() *domain.DBCChallengeInfo {
	if m == nil {
		return nil
	}

	obj := &domain.DBCChallengeInfo{
		Id:             int64(m.ID),
		OwnerId:        m.OwnerID,
		Name:           m.Name,
		Desc:           m.Desc,
		Image:          m.Image,
		IsAutoTrack:    m.IsAutoTrack,
		VisibilityType: m.VisibilityType,
		UpdatedAt:      m.UpdatedAt,
		CreatedAt:      m.CreatedAt,
		DeletedAt:      nil,
	}
	if m.DeletedAt.Valid {
		obj.DeletedAt = new(time.Time)
		*obj.DeletedAt = m.DeletedAt.Time
	}
	if m.Category != nil {
		obj.Category = m.Category.DTO()
	}

	return obj
}

type DBCChallengesUsers struct {
	gorm.Model

	UserID int64
	User   *User

	ChallengeID int64
	Challenge   *DBCChallenge

	LastSeries int64
}

func (m *DBCChallengesUsers) DTO() *domain.DBCUserChallenge {
	if m == nil {
		return nil
	}

	obj := &domain.DBCUserChallenge{
		Id:         int64(m.ID),
		UserId:     m.UserID,
		LastSeries: m.LastSeries,
		UpdatedAt:  m.UpdatedAt,
		CreatedAt:  m.CreatedAt,
		DeletedAt:  nil,
	}
	if m.DeletedAt.Valid {
		obj.DeletedAt = new(time.Time)
		*obj.DeletedAt = m.DeletedAt.Time
	}
	if m.Challenge != nil {
		obj.ChallengeInfo = m.Challenge.DTO()
	}

	return obj
}
