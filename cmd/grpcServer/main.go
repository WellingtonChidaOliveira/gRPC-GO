package main

import (
	"database/sql"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"gitub.com/wellingtonchida/gRPC/internal/database"
	"gitub.com/wellingtonchida/gRPC/internal/pb"
	"gitub.com/wellingtonchida/gRPC/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryService(*categoryDb)

	// Start gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(list); err != nil {
		panic(err)
	}
}
