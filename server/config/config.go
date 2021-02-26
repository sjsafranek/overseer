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
	Engine string
	Name   string
	Pass   string
	User   string
	Host   string
	Port   int64
}

func (self *Database) GetDatabaseConnection() string {
	return fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=disable",
		self.Engine,
		self.User,
		self.Pass,
		self.Host,
		self.Port,
		self.Name)
}
