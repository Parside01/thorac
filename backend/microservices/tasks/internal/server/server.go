package server

import (
	"context"
	"database/sql"
	"github.com/thorac/backend/tasks/internal/database"
	"github.com/thorac/backend/tasks/internal/proto/gen/go/tasks"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	LowPriority        = "low"
	MediumPriority     = "medium"
	HighPriority       = "high"
	ReallyHighPriority = "really_high"
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
	project, err := s.db.NewProject(ctx, database.NewProjectParams{
		Title:       req.Title,
		Description: sql.NullString{req.Description, true},
		UserID:      req.UserId,
	})

	if err != nil {
		s.logger.Error("Failed to create project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create project")
	}

	_, err = s.db.CreateNewTaskStatus(ctx, database.CreateNewTaskStatusParams{
		OwnerProjectID: project.ProjectID,
		Type:           "Ожидает",
	})
	if err != nil {
		s.logger.Error("Failed to create new task status", zap.Error(err))
	}

	_, err = s.db.CreateNewTaskStatus(ctx, database.CreateNewTaskStatusParams{
		OwnerProjectID: project.ProjectID,
		Type:           "В работе",
	})
	if err != nil {
		s.logger.Error("Failed to create new task status", zap.Error(err))
	}

	_, err = s.db.CreateNewTaskStatus(ctx, database.CreateNewTaskStatusParams{
		OwnerProjectID: project.ProjectID,
		Type:           "Завершена",
	})
	if err != nil {
		s.logger.Error("Failed to create new task status", zap.Error(err))
	}

	return &tasks.CreateProjectResponse{
		UserId:      project.UserID,
		Title:       project.Title,
		Description: project.Description.String,
		CreatedAt:   project.CreatedAt.Time.String(),
	}, nil
}

func (s *TaskServer) CreateTask(ctx context.Context, req *tasks.CreateTaskRequest) (*tasks.CreateTaskResponse, error) {
	task, err := s.db.NewTask(ctx, database.NewTaskParams{
		Title:       req.Title,
		Description: sql.NullString{req.Description, true},
		UserID:      req.UserId,
		Priority:    req.Priority,
	})
	if err != nil {
		s.logger.Error("Failed to create task", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create task")
	}
	return &tasks.CreateTaskResponse{
		UserId:      task.UserID,
		Title:       task.Title,
		Description: task.Description.String,
		Priority:    task.Priority,
	}, nil
}

func (s *TaskServer) CreateTaskStatus(ctx context.Context, req *tasks.CreateTaskStatusRequest) (*tasks.CreateTaskStatusRequest, error) {
	bool, err := s.db.IsUserInProject(ctx, database.IsUserInProjectParams{
		ProjectID: req.ProjectId,
		UserID:    req.UserId,
	})
	if err != nil {
		s.logger.Error("Failed to check user status in project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to check user status in project")
	}
	if !bool {
		return nil, status.Error(codes.FailedPrecondition, "User status in project doesn't exist")
	}

	stat, err := s.db.CreateNewTaskStatus(ctx, database.CreateNewTaskStatusParams{
		OwnerProjectID: req.ProjectId,
		Type:           req.Type,
	})
	if err != nil {
		s.logger.Error("Failed to create new task status", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create new task status")
	}
	return &tasks.CreateTaskStatusRequest{
		Type:      stat.Type,
		ProjectId: req.ProjectId,
		UserId:    req.UserId,
	}, nil
}

func (s *TaskServer) IsUserInProject(ctx context.Context, req *tasks.IsUserInProjectRequest) (*tasks.IsUserInProjectResponse, error) {
	in, err := s.db.IsUserInProject(ctx, database.IsUserInProjectParams{
		ProjectID: req.ProjectId,
		UserID:    req.UserId,
	})
	if err != nil {
		s.logger.Error("Failed to check user status in project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to check user status in project")
	}

	return &tasks.IsUserInProjectResponse{InProject: in}, nil
}

func (s *TaskServer) AddUserToProject(ctx context.Context, req *tasks.AddUserToProjectRequest) (*tasks.AddUserToProjectRequest, error) {
	_, err := s.db.AddUserToProject(ctx, database.AddUserToProjectParams{
		UserID:    req.UserId,
		ProjectID: req.ProjectId,
	})
	if err != nil {
		s.logger.Error("Failed to add user status to project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to add user status to project")
	}

	return &tasks.AddUserToProjectRequest{UserId: req.UserId}, nil
}

func (s *TaskServer) AddTaskToProject(ctx context.Context, req *tasks.AddTaskToProjectRequest) (*tasks.AddTaskToProjectRequest, error) {
	bool, err := s.db.IsUserInProject(ctx, database.IsUserInProjectParams{
		ProjectID: req.ProjectId,
		UserID:    req.UserId,
	})
	if err != nil {
		s.logger.Error("Failed to check user status in project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to check user status in project")
	}
	if !bool {
		return nil, status.Error(codes.FailedPrecondition, "User status in project doesn't exist")
	}

	err = s.db.AddTaskToProject(ctx, database.AddTaskToProjectParams{
		TaskID:    req.TaskId,
		ProjectID: req.ProjectId,
	})
	if err != nil {
		s.logger.Error("Failed to add task to project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to add task to project")
	}
	return &tasks.AddTaskToProjectRequest{TaskId: req.TaskId}, nil
}

func (s *TaskServer) UpdateProject(ctx context.Context, req *tasks.UpdateProjectRequest) (*tasks.UpdateProjectResponse, error) {
	proj, err := s.db.UpdateProject(ctx, database.UpdateProjectParams{
		Title:       req.Title,
		Description: sql.NullString{req.Description, true},
		UserID:      req.UserId,
	})
	if err != nil {
		s.logger.Error("Failed to update project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to update project")
	}
	return &tasks.UpdateProjectResponse{
		ProjectId:   proj.ProjectID,
		UserId:      proj.UserID,
		Title:       proj.Title,
		Description: proj.Description.String,
		CreatedAt:   proj.CreatedAt.Time.String(),
	}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *tasks.UpdateTaskRequest) (*tasks.UpdateTaskResponse, error) {
	has, err := s.db.UserHasThisTask(ctx, database.UserHasThisTaskParams{
		UserID: req.UserId,
		TaskID: req.TaskId,
	})
	if err != nil {
		s.logger.Error("Failed to check user status in project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to check user status in project")
	}

	if !has {
		return nil, status.Error(codes.FailedPrecondition, "User status in project doesn't exist")
	}

	t, err := s.db.UpdateTask(ctx, database.UpdateTaskParams{
		TaskID:      req.TaskId,
		Title:       req.Title,
		Description: sql.NullString{req.Description, true},
		UserID:      req.UserId,
		Priority:    req.Priority,
	})
	if err != nil {
		s.logger.Error("Failed to update task", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to update task")
	}

	return &tasks.UpdateTaskResponse{
		TaskId:      t.TaskID,
		Title:       t.Title,
		Description: t.Description.String,
		UserId:      t.UserID,
		Priority:    t.Priority,
	}, nil
}

func (s *TaskServer) AddUserToTask(ctx context.Context, req *tasks.AddUserToTaskRequest) (*tasks.Empty, error) {
	err := s.db.AddUserToTask(ctx, database.AddUserToTaskParams{
		UserID: req.UserId,
		TaskID: req.TaskId,
	})
	if err != nil {
		s.logger.Error("Failed to add user status to task", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to add user status to task")
	}

	return &tasks.Empty{}, nil
}

func (s *TaskServer) DeleteProject(ctx context.Context, req *tasks.DeleteProjectRequest) (*tasks.Empty, error) {
	if err := s.db.DeleteProject(ctx, req.ProjectId); err != nil {
		s.logger.Error("Failed to delete project", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to delete project")
	}
	if err := s.db.DeleteProjectsTask(ctx, req.ProjectId); err != nil {
		s.logger.Error("Failed to delete project task", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to delete project task")
	}
	if err := s.db.DeleteProjectFromUsersProjects(ctx, req.ProjectId); err != nil {
		s.logger.Error("Failed to delete project from users projects", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to delete project from users projects")
	}
	return &tasks.Empty{}, nil
}

func (s *TaskServer) GetUserProjects(ctx context.Context, req *tasks.GetUserProjectsRequest) (*tasks.GetUserProjectsResponse, error) {
	ids, err := s.db.GetUserProjects(ctx, req.UserId)
	if err != nil {
		s.logger.Error("Failed to get user projects", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get user projects")
	}
	var projs []*tasks.Project
	for _, id := range ids {
		proj, err := s.db.GetProjectByID(ctx, id.ProjectID)
		if err != nil {
			s.logger.Error("Failed to get project", zap.Error(err))
			return nil, status.Error(codes.Internal, "Failed to get project")
		}
		projs = append(projs, &tasks.Project{
			ProjectId:   proj.ProjectID,
			Title:       proj.Title,
			Description: proj.Description.String,
			CreatedAt:   proj.CreatedAt.Time.String(),
			UserId:      proj.UserID,
		})
	}

	return &tasks.GetUserProjectsResponse{
		Projects: projs,
	}, nil
}

func (s *TaskServer) GetUserTasks(ctx context.Context, req *tasks.GetUserTasksRequest) (*tasks.GetUserTasksResponse, error) {
	ids, err := s.db.GetUserTasks(ctx, req.UserId)
	if err != nil {
		s.logger.Error("Failed to get user tasks", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get user tasks")
	}
	var task []*tasks.Task
	for _, id := range ids {
		t, err := s.db.GetTaskById(ctx, id.TaskID)
		if err != nil {
			s.logger.Error("Failed to get task", zap.Error(err))
		}
		task = append(task, &tasks.Task{
			TaskId:      t.TaskID,
			UserId:      t.UserID,
			Title:       t.Title,
			Description: t.Description.String,
			Priority:    t.Priority,
		})
	}
	return &tasks.GetUserTasksResponse{
		Tasks: task,
	}, nil
}
