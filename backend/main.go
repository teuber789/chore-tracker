package main

import (
	"fmt"
	"log"
	"net"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/teuber789/chore-tracker/internal"
	"github.com/teuber789/chore-tracker/internal/gen"
	"google.golang.org/grpc"
)

func main() {
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

	// Start GRPC server
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
	server := internal.NewImplementedServer(db)
	gen.RegisterChoreTrackerServer(grpcServer, server)
	grpcServer.Serve(lis)
}
