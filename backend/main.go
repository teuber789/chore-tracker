package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/teuber789/chore-tracker/internal"
	"github.com/teuber789/chore-tracker/internal/gen"
	"google.golang.org/grpc"
)

// Starts and serves a GRPC server
func serveGrpc(db internal.ChoreTrackerStore) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 9090))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
		)),
	)
	server := internal.NewGrpcServer(db)
	gen.RegisterChoreTrackerServer(grpcServer, server)
	grpcServer.Serve(lis)
}

// Starts and serves an HTTP server
func serveHttp(db internal.ChoreTrackerStore) {
	r, err := internal.NewHttpRouter(db)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:8081",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Read in type of server to start
	serverType := flag.String("server", "", "Specifies whether to start the GRPC or the HTTP server. Valid values are 'grpc' and 'http'.")
	flag.Parse()

	// Run DB migrations
	// Source: https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md#optional-run-migrations-within-your-go-app
	m, err := migrate.New(
		"file://internal/db/migrations",
		internal.ConnString())
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	// Get DB handle
	db, err := internal.NewChoreTrackerStore()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// IRL, there would be no need to serve both an HTTP and a GRPC service from the same application.
	if *serverType == "grpc" {
		serveGrpc(db)
	} else if *serverType == "http" {
		serveHttp(db)
	} else {
		log.Fatalf("Unknown server type %s", *serverType)
	}
}
