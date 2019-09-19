package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/arkadyb/demos/blog2/internal/server"
	"github.com/arkadyb/demos/blog2/proto/reminder/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	serverCert, err := credentials.NewServerTLSFromFile("../server.crt", "../server.key")
	if err != nil {
		log.Fatalln("failed to create cert", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(serverCert))
	reminder.RegisterReminderServiceServer(grpcServer, new(server.MyServer))

	go func() {
		lis, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("failed to start server", err)
		}
	}()

	// let us wait for an input here (ctrl+c) to stop the client
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	signal := <-c
	log.Fatalf("process killed with signal: %v", signal.String())
}
