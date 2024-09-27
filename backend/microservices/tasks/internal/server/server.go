package server

import (
	"github.com/jiyeyuran/go-config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ServerConfig struct {
	Address string `json:"address" yaml:"address"`
}

type Server struct {
	logger      *zap.Logger
	router      *echo.Echo
	controllers []Controller
	Config      *ServerConfig
}

func NewServer(path string) *Server {
	serve := &Server{}
	serve.logger, _ = zap.NewDevelopment()
	if err := config.LoadFile(path); err != nil {
		serve.logger.Fatal("Error read config", zap.Error(err))
		return nil
	}
	if err := config.LoadFile(path); err != nil {
		serve.logger.Fatal("Error read config", zap.Error(err))
	}

	conf := &ServerConfig{}
	if err := config.Get("server").Scan(&conf); err != nil {
		serve.logger.Fatal("Error read config", zap.Error(err))
	}

	serve.router = echo.New()
	serve.Config = conf
	return serve
}

func (server *Server) Start() error {
	server.logger.Info("Starting server", zap.String("address", server.Config.Address))
	return server.router.Start(server.Config.Address)
}

func (s *Server) RegisterRouters(routes []Controller) {
	for _, route := range routes {
		group := s.router.Group(route.GetGroup())
		for _, handler := range route.GetHandlers() {
			group.Add(handler.GetMethod(), handler.GetPath(), handler.GetHandler())
		}
	}
	s.controllers = routes
}
