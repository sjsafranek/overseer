/*
 *
 * Copyright 2021 stefan safranek
 *
 */

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"

	pb "github.com/sjsafranek/overseer/service"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	"github.com/sjsafranek/logger"
	"github.com/sjsafranek/overseer/server/config"
	"github.com/sjsafranek/overseer/server/database"
)

const (
	PROJECT                   string = "Overseer"
	VERSION                   string = "0.0.3"
	DEFAULT_HOST              string = ""
	DEFAULT_PORT              int    = 50051
	DEFAULT_DATABASE_ENGINE   string = "postgres"
	DEFAULT_DATABASE_DATABASE string = "overseer"
	DEFAULT_DATABASE_PASSWORD string = "dev"
	DEFAULT_DATABASE_USERNAME string = "overseer"
	DEFAULT_DATABASE_HOST     string = "localhost"
	DEFAULT_DATABASE_PORT     int64  = 5432
	DEFAULT_REDIS_PORT        int64  = 6379
	DEFAULT_REDIS_HOST        string = ""
	DEFAULT_CONFIG_FILE       string = "config.json"
)

var (
	HOST              string = DEFAULT_HOST
	PORT              int    = DEFAULT_PORT
	DATABASE_ENGINE   string = DEFAULT_DATABASE_ENGINE
	DATABASE_DATABASE string = DEFAULT_DATABASE_DATABASE
	DATABASE_PASSWORD string = DEFAULT_DATABASE_PASSWORD
	DATABASE_USERNAME string = DEFAULT_DATABASE_USERNAME
	DATABASE_HOST     string = DEFAULT_DATABASE_HOST
	DATABASE_PORT     int64  = DEFAULT_DATABASE_PORT
	REDIS_PORT        int64  = DEFAULT_REDIS_PORT
	REDIS_HOST        string = DEFAULT_REDIS_HOST
	CONFIG_FILE       string = DEFAULT_CONFIG_FILE
	db                *database.Database
	conf              *config.Config
)

type Server struct {
	pb.UnimplementedOverseerServer
}

func successResponse(user *database.User, apikey *database.Apikey) (*pb.Response, error) {
	response := pb.Response{Status: "ok"}

	if nil != user {

		var apikeys []*pb.Apikey
		for _, apikey := range user.Apikeys {
			apikeys = append(apikeys, &pb.Apikey{
				UserId:    apikey.UserId,
				Name:      apikey.Name,
				Apikey:    apikey.Apikey,
				Secret:    apikey.Secret,
				IsActive:  apikey.IsActive,
				IsDeleted: apikey.IsDeleted,
				CreatedAt: apikey.CreatedAt.String(),
				UpdatedAt: apikey.UpdatedAt.String(),
			})
		}

		response.User = &pb.User{
			Id:        user.Id,
			Username:  user.Username,
			Email:     user.Email,
			IsActive:  user.IsActive,
			IsDeleted: user.IsDeleted,
			CreatedAt: user.CreatedAt.String(),
			UpdatedAt: user.UpdatedAt.String(),
			Apikeys:   apikeys,
		}
	}

	if nil != apikey {
		response.Apikey = &pb.Apikey{
			UserId:    apikey.UserId,
			Name:      apikey.Name,
			Apikey:    apikey.Apikey,
			Secret:    apikey.Secret,
			IsActive:  apikey.IsActive,
			IsDeleted: apikey.IsDeleted,
			CreatedAt: apikey.CreatedAt.String(),
			UpdatedAt: apikey.UpdatedAt.String(),
		}
	}

	return &response, nil
}

// CreateUser creates new user
func (self *Server) CreateUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	user, err := db.CreateUserIfNotExists(req.GetEmail(), req.GetUsername())
	if nil != err {
		return &pb.Response{}, err
	}
	return successResponse(user, nil)
}

// getUserByUsername retrieves user from database and passes it to callback function
func getUserByUsername(username string, clbk func(*database.User) (*pb.Response, error)) (*pb.Response, error) {
	user, err := db.GetUserByUsername(username)
	if nil != err {
		return &pb.Response{}, err
	}
	return clbk(user)
}

// AuthenticateUser checks user password
func (self *Server) AuthenticateUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		ok, err := user.IsPassword(req.GetPassword())
		if nil != err {
			return &pb.Response{}, err
		}
		if !ok {
			return &pb.Response{}, errors.New("Incorrect password")
		}
		return successResponse(user, nil)
	})

}

// GetUser from database
func (self *Server) GetUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		return successResponse(user, nil)
	})
}

// DeleteUser markers user has deleted
func (self *Server) DeleteUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		err := user.Delete()
		if nil != err {
			return &pb.Response{}, err
		}
		return successResponse(nil, nil)
	})
}

// ActivateUser marks user as active
func (self *Server) ActivateUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		err := user.Activate()
		if nil != err {
			return &pb.Response{}, err
		}
		return successResponse(nil, nil)
	})
}

// DeactivateUser marks user as deactive
func (self *Server) DeactivateUser(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		err := user.Deactivate()
		if nil != err {
			return &pb.Response{}, err
		}
		return successResponse(nil, nil)
	})
}

// SetPassword sets new password for user
func (self *Server) SetUserPassword(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		err := user.SetPassword(req.GetPassword())
		if nil != err {
			return &pb.Response{}, err
		}
		return successResponse(nil, nil)
	})
}

// CreateUserApikey
func (self *Server) CreateUserApikey(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		apikey, err := user.CreateApikey(req.GetName())
		if nil != err {
			return &pb.Response{}, err
		}
		return successResponse(nil, apikey)
	})
}

// CreateUserSocialAccount
func (self *Server) CreateUserSocialAccount(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return getUserByUsername(req.GetUsername(), func(user *database.User) (*pb.Response, error) {
		err := user.CreateSocialAccountIfNotExists(req.GetId(), req.GetName(), req.GetType())
		if nil != err {
			return &pb.Response{}, err
		}
		return successResponse(nil, nil)
	})
}

func init() {
	var printVersion bool

	// read credentials from environment variables if available
	conf = &config.Config{
		Server: config.Server{
			Host: DEFAULT_HOST,
			Port: DEFAULT_PORT,
		},
		Database: config.Database{
			DatabaseEngine: DATABASE_ENGINE,
			DatabaseHost:   DEFAULT_DATABASE_HOST,
			DatabaseName:   DEFAULT_DATABASE_DATABASE,
			DatabasePass:   DEFAULT_DATABASE_PASSWORD,
			DatabaseUser:   DEFAULT_DATABASE_USERNAME,
			DatabasePort:   DEFAULT_DATABASE_PORT,
		},
		Redis: config.Redis{
			Host: DEFAULT_REDIS_HOST,
			Port: DEFAULT_REDIS_PORT,
		},
	}

	flag.StringVar(&conf.Server.Host, "Host", DEFAULT_HOST, "Server host")
	flag.IntVar(&conf.Server.Port, "port", DEFAULT_PORT, "Server port")
	flag.StringVar(&conf.Database.DatabaseHost, "dbhost", DEFAULT_DATABASE_HOST, "database host")
	flag.StringVar(&conf.Database.DatabaseName, "dbname", DEFAULT_DATABASE_DATABASE, "database name")
	flag.StringVar(&conf.Database.DatabasePass, "dbpass", DEFAULT_DATABASE_PASSWORD, "database password")
	flag.StringVar(&conf.Database.DatabaseUser, "dbuser", DEFAULT_DATABASE_USERNAME, "database username")
	flag.Int64Var(&conf.Database.DatabasePort, "dbport", DEFAULT_DATABASE_PORT, "Database port")
	flag.StringVar(&conf.Redis.Host, "redishost", DEFAULT_REDIS_HOST, "Redis host")
	flag.Int64Var(&conf.Redis.Port, "redisport", DEFAULT_REDIS_PORT, "Redis port")
	flag.StringVar(&CONFIG_FILE, "c", DEFAULT_CONFIG_FILE, "config file")
	flag.BoolVar(&printVersion, "V", false, "Print version and exit")
	flag.Parse()

	if printVersion {
		fmt.Println(PROJECT, VERSION)
		os.Exit(0)
	}
}

func main() {

	logger.Infof("%v v%v", PROJECT, VERSION)
	logger.Debug("GOOS: ", runtime.GOOS)
	logger.Debug("CPUS: ", runtime.NumCPU())
	logger.Debug("PID: ", os.Getpid())
	logger.Debug("Go Version: ", runtime.Version())
	logger.Debug("Go Arch: ", runtime.GOARCH)
	logger.Debug("Go Compiler: ", runtime.Compiler)
	logger.Debug("NumGoroutine: ", runtime.NumGoroutine())

	db = database.New(conf.Database.GetDatabaseConnection())

	version, err := db.GetVersion()
	if nil != err {
		panic(err)
	}
	logger.Debugf("Database version: %v", version)

	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", HOST, PORT))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterOverseerServer(server, &Server{})
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
