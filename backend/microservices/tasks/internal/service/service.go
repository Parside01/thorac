package service

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thorac/backend/tasks/internal/proto/gen/go/tasks"
	"github.com/thorac/backend/tasks/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Panicf("Error load .env file %s", err.Error())
	}
}

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panicf("Error create logger")
	}

	conn, err := net.Listen("tcp", os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error listening on address %s", os.Getenv("SERVER_ADDRESS")), zap.Error(err))
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_CONN"))
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error connecting to database on address %s", os.Getenv("POSTGRES_CONN")), zap.Error(err))
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		logger.Fatal(fmt.Sprintf("Error pinging database on address %s", os.Getenv("POSTGRES_CONN")), zap.Error(err))
	}
	logger.Info(fmt.Sprintf("Connect to database on %s address", os.Getenv("POSTGRES_CONN")))

	task_server := server.NewTaskServer(logger, db)
	grpc_server := grpc.NewServer()

	tasks.RegisterTasksServer(grpc_server, task_server)

	if err := grpc_server.Serve(conn); err != nil {
		logger.Fatal(fmt.Sprintf("Error starting grpc server on address %s", os.Getenv("SERVER_ADDRESS")), zap.Error(err))
	}

	logger.Info(fmt.Sprintf("Listening on address %s", os.Getenv("SERVER_ADDRESS")))
}
