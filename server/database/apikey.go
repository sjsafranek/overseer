package database

import (
	"time"
)

// Apikey struct of apikey object
type Apikey struct {
	UserId    string    `json:"user_id"`
	Name      string    `json:"name"`
	Apikey    string    `json:"apikey"`
	Secret    string    `json:"secret"`
	IsDeleted bool      `json:"is_deleted"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at,string"`
	UpdatedAt time.Time `json:"updated_at,string"`
	db        *Database `json:"-"`
}

// SetDatabase set object database reference
func (self *Apikey) SetDatabase(db *Database) {
	self.db = db
}

// Delete markers object as deleted
func (self *Apikey) Delete() error {
	self.IsDeleted = true
	return self.Update()
}

// Activate markers object as active
func (self *Apikey) Activate() error {
	self.IsActive = true
	return self.Update()
}

// Deactivate markers object as inactive
func (self *Apikey) Deactivate() error {
	self.IsActive = false
	return self.Update()
}

// Update updates object data in database
func (self *Apikey) Update() error {
	return self.db.Insert(`
		UPDATE apikeys
			SET
				is_deleted=$1,
				is_active=$2
			WHERE apikey=$3;`, self.IsDeleted, self.IsActive, self.Apikey)
}

//
