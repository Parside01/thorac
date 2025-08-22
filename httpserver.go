package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type HttpServer struct {
	config *Config
	node   Node
	logger *zap.Logger
	router *echo.Echo
}

func NewHttpServer(config *Config, logger *zap.Logger, node Node) *HttpServer {
	router := echo.New()
	router.Use(middleware.Recover())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger())

	server := &HttpServer{
		config: config,
		logger: logger,
		router: router,
		node:   node,
	}

	server.RegisterRoutes()
	return server
}

func (s *HttpServer) RegisterRoutes() {
	s.router.GET("/health", s.handleHealth)
	s.router.POST("/join", s.handleJoin)
}

func (s *HttpServer) Start() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		_ = s.Shutdown()
	}()

	s.logger.Info("starting http server", zap.String("server", s.config.ServerAddr))
	return s.router.Start(s.config.ServerAddr)
}

func (s *HttpServer) handleHealth(c echo.Context) error {
	return c.JSON(http.StatusOK, s.node.Health())
}

func (s *HttpServer) handleJoin(c echo.Context) error {
	id := c.QueryParam("id")
	addr := c.QueryParam("addr")

	if id == "" || addr == "" {
		s.logger.Warn("invalid join request: missing id or addr")
		return echo.NewHTTPError(http.StatusBadRequest, "query params 'id' and 'addr' are required")
	}

	s.logger.Info("attempting to join node", zap.String("peer_id", id), zap.String("peer_address", addr))

	if err := s.node.Join(id, addr); err != nil {
		s.logger.Error("join failed", zap.String("peer_id", id), zap.String("peer_address", addr), zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, "success")
}

func (s *HttpServer) Shutdown() error {
	return s.router.Shutdown(context.Background())
}
