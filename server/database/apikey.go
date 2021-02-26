package database

import (
	"time"
)

type Apikey struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Apikey    string    `json:"apikey"`
	Secret    string    `json:"secret"`
	IsDeleted bool      `json:"is_deleted"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at,string"`
	UpdatedAt time.Time `json:"updated_at,string"`
	db          *Database `json:"-"`
}

func (self *Apikey) SetDatabase(db *Database) {
    self.db = db
}
