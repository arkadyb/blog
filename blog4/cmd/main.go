package main

import (
	"context"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/arkadyb/blog/blog4/internal/server"
	"github.com/arkadyb/blog/blog4/proto/reminder/v1"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	serverCert, err := credentials.NewServerTLSFromFile("../server.crt", "../server.key")
	if err != nil {
		log.Fatalln("failed to create cert", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_validator.UnaryServerInterceptor(),
	)), grpc.Creds(serverCert))
	reminder.RegisterReminderServiceServer(grpcServer, new(server.MyServer))

	clientCert, err := credentials.NewClientTLSFromFile("../server.crt", "")
	if err != nil {
		log.Fatalln("failed to create cert", err)
	}
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8080",
		grpc.WithTransportCredentials(clientCert),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	router := runtime.NewServeMux()
	if err = reminder.RegisterReminderServiceHandler(context.Background(), router, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	http.ListenAndServeTLS(":8080", "../server.crt", "../server.key", httpGrpcRouter(grpcServer, router))
}

func httpGrpcRouter(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}
