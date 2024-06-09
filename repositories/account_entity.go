package repositories

import (
	"encoding/json"
	"github.com/google/uuid"
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

type Stakeholder struct {
	bun.BaseModel  `bun:"table:stakeholders"`
	ID             int64           `bun:"id,pk,autoincrement"       json:"-"`
	UID            uuid.UUID       `json:"uid,omitempty"`
	Properties     json.RawMessage `json:"properties,omitempty"`
	Roles          []string        `bun:",array"                    json:"roles,omitempty"`
	OrganizationID int             `json:"organizationId,omitempty"`
	CreatedAt      time.Time       `json:"-"`
	UpdatedAt      *time.Time      `json:"-"`
	DeletedAt      *time.Time      `bun:",soft_delete"              json:"-"`
}

type AccountEntityCreateParams struct {
	bun.BaseModel `bun:"table:accounts"`
	Name          string `json:"name"`
}

type AccountEntityUpdateParams struct {
	bun.BaseModel `bun:"table:accounts"`
	Name          string `json:"name"`
}
