package server

import (
	"database/sql"
	"github.com/jiyeyuran/go-config"
	"github.com/labstack/echo/v4"
	"github.com/thorac/backend/tasks/internal/database"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type TaskControllerConfig struct {
	RouteGroup string `yaml:"route_group"`
}

type TaskController struct {
	logger *zap.Logger
	config *TaskControllerConfig
	db     *database.Queries
}

func NewTaskController(logger *zap.Logger, db *sql.DB) *TaskController {
	conf := &TaskControllerConfig{}
	if err := config.Get("task_controller").Scan(&conf); err != nil {
		logger.Fatal("Error read task_controller config", zap.Error(err))
	}
	return &TaskController{
		logger: logger,
		config: conf,
		db:     database.New(db),
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

type CreateProjectRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name"`
}

func (c *TaskController) createProject(e echo.Context) error {
	userId, err := strconv.ParseInt(e.Request().Header.Get("UserID"), 10, 64)
	if err != nil {
		c.logger.Error("NotAuthorization user")
		return echo.NewHTTPError(http.StatusUnauthorized, "NotAuthorization user")
	}

	var req CreateProjectRequest
	if err := e.Bind(&req); err != nil {
		c.logger.Error("Bad request", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}

	par := database.CreateProjectParams{
		Title:       req.Title,
		Description: sql.NullString{String: req.Description, Valid: true},
		AuthorName:  req.AuthorName,
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
		UserID:      userId,
	}
	proj, err := c.db.CreateProject(e.Request().Context(), par)
	if err != nil {
		c.logger.Error("Error creating project", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "Error creating project")
	}
	return echo.NewHTTPError(http.StatusCreated, proj)
}
