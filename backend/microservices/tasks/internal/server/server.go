package server

import (
	"context"
	"database/sql"
	"github.com/thorac/backend/tasks/internal/database"
	"github.com/thorac/backend/tasks/internal/proto/gen/go/tasks"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type TaskServer struct {
	logger *zap.Logger
	db     *database.Queries
	tasks.UnimplementedTasksServer
}

func NewTaskServer(logger *zap.Logger, db *sql.DB) *TaskServer {
	return &TaskServer{
		logger: logger,
		db:     database.New(db),
	}
}

func (s *TaskServer) CreateProject(ctx context.Context, req *tasks.CreateProjectRequest) (*tasks.CreateProjectResponse, error) {
	project, err := s.db.CreateProject(ctx, database.CreateProjectParams{
		Title:       req.Title,
		Description: sql.NullString{req.Description, true},
		UserID:      req.UserId,
	})

	if err != nil {
		s.logger.Error("Failed to create project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create project")
	}

	_, err = s.db.AddUserToProject(ctx, database.AddUserToProjectParams{
		ProjectID: project.ProjectID,
		UserID:    project.UserID,
	})
	if err != nil {
		s.logger.Error("Failed to add user to project", zap.Error(err))
	}

	return &tasks.CreateProjectResponse{
		UserId:      project.UserID,
		Title:       project.Title,
		Description: project.Description.String,
		CreatedAt:   project.CreatedAt.Time.String(),
	}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *tasks.CreateTaskRequest) (*tasks.CreateTaskResponse, error) {
	deadline, _ := time.Parse(time.Layout, req.Deadline)

	task, err := s.db.CreateTask(ctx, database.CreateTaskParams{
		Title: req.Title,
		Description: sql.NullString{req.Description, true},
		UserID: req.UserId,
		Deadline: sql.NullTime{time.Parse("", req.Deadline), true}
	})

	return &tasks.CreateTaskResponse{}, nil
}

func (s *TaskServer) CreateNewTaskStatus(ctx context.Context, req *tasks.CreateNewTaskStatusRequest) (*tasks.CreateNewTaskStatusResponse, error) {
	return &tasks.CreateNewTaskStatusResponse{}, nil
}

func (s *TaskServer) AddTaskToProject(ctx context.Context, req *tasks.AddTaskToProjectRequest) (*tasks.AddTaskToProjectResponse, error) {
	// TO DO: implement AddTaskToProject logic
	return &tasks.AddTaskToProjectResponse{}, nil
}

func (s *TaskServer) SetTaskStatus(ctx context.Context, req *tasks.SetTaskStatusRequest) (*tasks.SetTaskStatusResponse, error) {
	// TO DO: implement SetTaskStatus logic
	return &tasks.SetTaskStatusResponse{}, nil
}

func (s *TaskServer) DeleteTaskById(ctx context.Context, req *tasks.DeleteTaskByIdRequest) (*tasks.DeleteTaskByIdResponse, error) {
	// TO DO: implement DeleteTaskById logic
	return &tasks.DeleteTaskByIdResponse{}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *tasks.UpdateTaskRequest) (*tasks.UpdateTaskResponse, error) {
	// TO DO: implement UpdateTask logic
	return &tasks.UpdateTaskResponse{}, nil
}

func (s *TaskServer) GetTaskById(context.Context, *tasks.GetProjectTasksRequest) (*tasks.Task, error) {
	// TO DO: implement GetTaskById logic
	return &tasks.Task{}, nil
}

func (s *TaskServer) GetProjectTasks(ctx context.Context, req *tasks.GetProjectTasksRequest) (*tasks.TaskList, error) {
	// TO DO: implement GetProjectTasks logic
	return &tasks.TaskList{}, nil
}

func (s *TaskServer) GetTaskStatusById(ctx context.Context, req *tasks.GetTaskStatusByIdRequest) (*tasks.GetTaskStatusByIdResponse, error) {
	// TO DO: implement GetTaskStatusById logic
	return &tasks.GetTaskStatusByIdResponse{}, nil
}

func (s *TaskServer) GetTasksByStatus(ctx context.Context, req *tasks.GetTasksByStatusRequest) (*tasks.TaskList, error) {
	// TO DO: implement GetTasksByStatus logic
	return &tasks.TaskList{}, nil
}

func (s *TaskServer) GetTasksByPriority(ctx context.Context, req *tasks.GetTasksByPriorityRequest) (*tasks.TaskList, error) {
	// TO DO: implement GetTasksByPriority logic
	return &tasks.TaskList{}, nil
}

func (s *TaskServer) GetTasksByUserId(ctx context.Context, req *tasks.GetTasksByUserIdRequest) (*tasks.TaskList, error) {
	// TO DO: implement GetTasksByUserId logic
	return &tasks.TaskList{}, nil
}
