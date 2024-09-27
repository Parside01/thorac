package server

import (
	"github.com/jiyeyuran/go-config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

type Server struct {
	Logger *zap.Logger
	Config *config.Config
	Router *http.ServeMux
}

func NewServer(conf *config.Config) *Server {
	serve := &Server{}
	configurateLogger(conf)
	return serve
}

func configurateLogger(conf *config.Config) zapcore.Core {
	return zapcore.NewNopCore()
}
