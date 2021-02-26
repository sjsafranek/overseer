/*
 *
 * Copyright 2021 stefan safranek
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"flag"

	"google.golang.org/grpc"
	pb "github.com/sjsafranek/goauthserver/service"
)

func main() {

	var username string
	var password string
	var host string
	var port int64
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
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
	c := pb.NewAuthenticationClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AuthenticateUser(ctx, &pb.Request{Username: username, Password: password})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", r.GetUser().GetUsername())
}
