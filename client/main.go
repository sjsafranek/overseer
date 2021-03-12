/*
 *
 * Copyright 2021 stefan safranek
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/sjsafranek/overseer/service"
	"google.golang.org/grpc"
)

func messageToJSON(message proto.Message) string {
	jbytes, _ := protojson.Marshal(message)
	return string(jbytes)
}


func main() {
	var name string
	var username string
	var password string
	var email string
	var host string
	var port int64
	flag.StringVar(&name, "name", "", "name")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&email, "email", "", "email")
	flag.StringVar(&host, "host", "localhost", "server host")
	flag.Int64Var(&port, "port", 50051, "server port")
	flag.Parse()

	address := fmt.Sprintf("%v:%v", host, port)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewOverseerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Get action command from command line args
	action := "login"
	args := flag.Args()
	if 0 != len(args) {
		action = args[0]
	}

	if "create_user" == action {
		user, err := client.CreateUser(ctx, &pb.Request{Username: username, Email: email})
		if err != nil {
			log.Fatal(err)
		}
		if "" != password {
			_, err = client.SetUserPassword(ctx, &pb.Request{Username: username, Password: password})
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Println(messageToJSON(user))
		return
	}

	if "login" == action {
		response, err := client.AuthenticateUser(ctx, &pb.Request{Username: username, Password: password})
		if err != nil {
			log.Fatal(err)
		}
		user := response.GetUser()
		fmt.Println(messageToJSON(user))
		return
	}

	if "create_key" == action {
		response, err := client.CreateUserApikey(ctx, &pb.Request{Username: username, Name: name})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(messageToJSON(response.GetApikey()))
		return
	}
}
