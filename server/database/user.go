package database

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

// User struct of user object
type User struct {
	Id             string         `json:"id"`
	Username       string         `json:"username"`
	Password       string         `json:"-"`
	Email          string         `json:"email"`
	IsDeleted      bool           `json:"is_deleted"`
	IsActive       bool           `json:"is_active"`
	CreatedAt      time.Time      `json:"created_at,string"`
	UpdatedAt      time.Time      `json:"updated_at,string"`
	Apikeys        []Apikey       `json:"apikeys"`
	SocialAccounts *[]interface{} `json:"social_accounts"`
	db             *Database      `json:"-"`
}

// SetEmail sets user email
func (self *User) SetEmail(email string) error {
	self.Email = email
	return self.Update()
}

// SetPassword sets user password
func (self *User) SetPassword(password string) error {
	self.Password = password
	return self.Update()
}

// Delete markers object as deleted
func (self *User) Delete() error {
	self.IsDeleted = true
	return self.Update()
}

// Activate markers object as active
func (self *User) Activate() error {
	self.IsActive = true
	return self.Update()
}

// Deactivate markers object as inactive
func (self *User) Deactivate() error {
	self.IsActive = false
	return self.Update()
}

// Update updates object data in database
func (self *User) Update() error {
	return self.db.Insert(`
		UPDATE users
			SET
				email=$1,
				password=$2,
				is_deleted=$3,
				is_active=$4
			WHERE username=$5;`, self.Email, self.Password, self.IsDeleted, self.IsActive, self.Username)
}

// SetDatabase set object database reference
func (self *User) SetDatabase(db *Database) {
	self.db = db
}

// IsPassword checks if provided password/hash matches database record
func (self *User) IsPassword(password string) (bool, error) {
	match := false
	return match, self.db.Exec(func(conn *sql.DB) error {
		rows, err := conn.Query(`SELECT is_password($1, $2);`, self.Username, password)

		if nil != err {
			return err
		}

		for rows.Next() {
			rows.Scan(&match)
			return nil
		}

		return errors.New("Not found")
	})
}

// CreateApikey creates new apikey for user
func (self *User) CreateApikey(name string) (*Apikey, error) {
	var apikey Apikey
	return &apikey, self.db.Select(&apikey, `
		INSERT INTO apikeys (user_id, name) VALUES ($1, $2) RETURNING json_build_object(
			'user_id', user_id,
			'name', name,
			'apikey', apikey,
			'secret', secret,
			'is_active', is_active,
			'is_deleted', is_deleted,
			'created_at', to_char(created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
			'updated_at', to_char(updated_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"')
		);
	`, self.Id, name)
}

/**
 * Social Accounts
 */
// CreateSocialAccountIfNotExists
// https://stackoverflow.com/questions/4069718/postgres-insert-if-does-not-exist-already
// ON CONFLICT DO NOTHING/UPDATE
// http://www.postgresqltutorial.com/postgresql-upsert/
func (self *User) CreateSocialAccountIfNotExists(social_id, social_name, social_type string) error {
	err := self.db.Insert(`
		INSERT INTO social_accounts(id, name, type, user_id)
			VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING;
	`, social_id, social_name, social_type, self.Id)
	if nil != err && strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return nil
	}
	return nil
}
