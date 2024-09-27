package server

import (
	"github.com/jiyeyuran/go-config"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
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
	return []ControllerHandler{
		&Handler{
			Method:  http.MethodPost,
			Path:    "/create_project",
			Handler: c.createProject,
		},
	}
}

func (c *TaskController) createProject(e echo.Context) error {

}
