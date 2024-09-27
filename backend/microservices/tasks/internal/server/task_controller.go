package server

import (
	"github.com/jiyeyuran/go-config"
	"go.uber.org/zap"
)

type TaskControllerConfig struct {
	RouteGroup string `yaml:"route_group"`
}

type TaskController struct {
	logger *zap.Logger
	config *TaskControllerConfig
}

func NewTaskController(logger *zap.Logger) *TaskController {
	conf := &TaskControllerConfig{}
	if err := config.Get("task_controller").Scan(&conf); err != nil {
		logger.Fatal("Error read task_controller config", zap.Error(err))
	}
	return &TaskController{
		logger: logger,
		config: conf,
	}
}

func (c *TaskController) GetGroup() string {
	return c.config.RouteGroup
}

func (c *TaskController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{}
}
