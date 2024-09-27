package service

import (
	"github.com/jiyeyuran/go-config"
	"github.com/thorac/backend/tasks/internal/server"
	"log"
)

func init() {
	if err := config.LoadFile("config.json"); err != nil {
		log.Panicf("Error read config %s", err.Error())
	}
}

func Run() {
	serve := server.NewServer()
	serve.RegisterRouters([]server.Controller{
		server.NewTaskController(serve.Logger),
	})
	if err := serve.Start(); err != nil {
		serve.Logger.Fatal(err.Error())
	}
}
