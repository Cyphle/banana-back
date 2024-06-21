package repositories

import (
	"github.com/uptrace/bun"
	"time"
)

type AccountEntity struct {
	bun.BaseModel `bun:"table:accounts"`
	ID            int64      `bun:"id,pk,autoincrement"       json:"-"`
	Name          string     `json:"name,omitempty"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	DeletedAt     *time.Time `bun:",soft_delete"              json:"-"`
}

type AccountEntityCreateParams struct {
	bun.BaseModel `bun:"table:accounts"`
	Name          string `json:"name"`
}
