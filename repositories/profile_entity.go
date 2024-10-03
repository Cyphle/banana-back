package repositories

import (
	"github.com/uptrace/bun"
	"time"
)

type ProfileEntity struct {
	bun.BaseModel `bun:"table:profiles"`
	ID            int64      `bun:"id,pk,autoincrement"       json:"-"`
	Username      string     `json:"username,omitempty"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	DeletedAt     *time.Time `bun:",soft_delete"              json:"-"`
}
