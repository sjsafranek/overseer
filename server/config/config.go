package config

import (
	"fmt"
)

type Config struct {
	Server   Server
	Database Database
	Redis    Redis
}

type Server struct {
	Host string
	Port int
}

type Redis struct {
	Host string
	Port int64
}

func (self *Redis) GetConnectionString() string {
	return fmt.Sprintf("%v:%v", self.Host, self.Port)
}

type Database struct {
	DatabaseEngine string
	DatabaseName   string
	DatabasePass   string
	DatabaseUser   string
	DatabaseHost   string
	DatabasePort   int64
}

func (self *Database) GetDatabaseConnection() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
		self.DatabaseEngine,
		self.DatabaseUser,
		self.DatabasePass,
		self.DatabaseHost,
		self.DatabasePort,
		self.DatabaseName)
}
