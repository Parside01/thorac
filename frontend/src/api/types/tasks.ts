export interface CreateProjectRequest {
  user_id: string;
  title: string;
  description: string;
}

export interface CreateProjectResponse {
  user_id: string;
  project_id: string;
  title: string;
  description: string;
  created_at: string;
}

export interface CreateTaskRequest {
  user_id: string;
  title: string;
  description: string;
  priority: string;
  deadline: string;
}

export interface CreateTaskResponse {
  user_id: string;
  task_id: string;
  title: string;
  description: string;
  priority: string;
  created_at: string;
  deadline: string;
}

export interface AddTaskToProjectRequest {
  user_id: string;
  project_id: string;
  task_id: string;
}

export interface AddTaskToProjectResponse {
}

export interface CreateNewTaskStatusRequest {
  user_id: string;
  project_id: string;
  status_type: string;
}

export interface CreateNewTaskStatusResponse {
  status_id: string;
  owner_project_id: string;
  status_type: string;
}

export interface SetTaskStatusRequest {
  user_id: string;
  task_id: string;
  status_id: string;
}

export interface SetTaskStatusResponse {
}

export interface Task {
  user_id: string;
  task_id: string;
  title: string;
  description: string;
  priority: string;
  created_at: string;
  deadline: string;
}

export interface TaskList {
  tasks: Task[];
}

export interface GetProjectTasksRequest {
  user_id: string;
  project_id: string;
}

export interface GetTaskByIdRequest {
  user_id: string;
  task_id: string;
}

export interface DeleteTaskByIdRequest {
  user_id: string;
  task_id: string;
}

export interface DeleteTaskByIdResponse {
}

export interface UpdateTaskRequest {
  user_id: string;
  task_id: string;
  title: string;
  description: string;
  priority: string;
  deadline: string;
}

export interface UpdateTaskResponse {
}

export interface GetTaskStatusByIdRequest {
  user_id: string;
  task_id: string;
}

export interface GetTaskStatusByIdResponse {
  status_type: string;
}

export interface GetTasksByStatusRequest {
  user_id: string;
  project_id: string;
  status_id: string;
}

export interface GetTasksByPriorityRequest {
  user_id: string;
  project_id: string;
  priority: string;
}

export interface GetTasksByUserIdRequest {
  user_id: string;
}

export interface GetTasksInProjectRequest {
  user_id: string;
  project_id: string;
}
