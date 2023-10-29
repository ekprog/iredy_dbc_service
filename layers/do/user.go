package do

import (
	"gorm.io/gorm"
	"microservice/layers/domain"
	"time"
)

type User struct {
	gorm.Model
	Score int64
}

func (m *User) DTO() *domain.User {
	u := &domain.User{
		Id:        int64(m.ID),
		Score:     m.Score,
		UpdatedAt: m.UpdatedAt,
		CreatedAt: m.CreatedAt,
		DeletedAt: nil,
	}

	if m.DeletedAt.Valid {
		u.DeletedAt = new(time.Time)
		*u.DeletedAt = m.DeletedAt.Time
	}

	return u
}
