-- name: NewProject :one
INSERT INTO projects (title, description, user_id, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: NewTask :one
INSERT INTO tasks (title, description, priority, user_id, created_at)
VALUES ($1, $2, $3, $4, NOW())
RETURNING *;

-- name: IsUserInProject :one
SELECT EXISTS (
    SELECT 1
    FROM users_projects
    WHERE project_id = $1 AND user_id = $2
);

-- name: DeleteProjectFromUsersProjects :exec
DELETE FROM users_projects
WHERE project_id = $1;

-- name: DeleteProjectsTask :exec
DELETE FROM project_tasks
WHERE project_id = $1;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = $1;


-- name: AddUserToProject :one
INSERT INTO users_projects (user_id, project_id)
VALUES ($1, $2)
RETURNING *;


-- name: AddTaskToProject :exec
INSERT INTO project_tasks (task_id, project_id)
VALUES ($1, $2);

-- name: UpdateProject :one
UPDATE projects
SET title = $1, description = $2
FROM users_projects
WHERE projects.project_id = users_projects.project_id AND users_projects.user_id = $3
RETURNING *;

-- name: UserHasThisTask :one
SELECT COUNT(*) > 0 FROM users_tasks
WHERE user_id = $1 AND task_id = $2;

-- name: AddUserToTask :exec
INSERT INTO users_tasks (task_id, user_id)
VALUES ($1, $2)
ON CONFLICT (task_id, user_id) DO NOTHING;


-- name: UpdateTask :one
UPDATE tasks
SET title       = $2,
    description = $3,
    priority    = $4,
    user_id     = $5
WHERE task_id = $1
RETURNING *;

-- name: CreateNewTaskStatus :one
INSERT INTO task_status_types (owner_project_id, type)
VALUES ($1, $2)
RETURNING *;

-- name: SetTaskStatus :exec
UPDATE task_statuses ts
SET task_status_id = (SELECT tst.task_status_types_id
                      FROM task_status_types tst
                      WHERE tst.type = $2
                        AND tst.owner_project_id = (SELECT pt.project_id
                                                    FROM project_tasks pt
                                                    WHERE pt.task_id = $1))
WHERE ts.task_id = $1;

-- name: GetAllTasks :many
SELECT *
FROM tasks;

-- name: GetUserTasks :many
SELECT * FROM users_tasks
WHERE user_id = $1;

-- name: GetUserProjects :many
SELECT * FROM users_projects
WHERE user_id = $1;

-- name: GetTaskById :one
SELECT *
FROM tasks
WHERE task_id = $1;

-- name: DeleteTask :exec
DELETE
FROM tasks
WHERE task_id = $1;


-- name: GetProjectByID :one
SELECT * FROM projects
WHERE project_id = $1 LIMIT 1;

-- name: GetTaskStatusByTaskId :one
SELECT t.type
FROM task_statuses AS ts
         JOIN task_status_types AS t ON ts.task_status_id = t.task_status_types_id
WHERE ts.task_id = $1;

-- name: GetAllTaskStatuses :many
SELECT *
FROM task_status_types;

-- name: GetTasksByStatus :many
SELECT t.*
FROM tasks AS t
         JOIN task_statuses AS ts ON t.task_id = ts.task_id
         JOIN task_status_types AS tst ON ts.task_status_id = tst.task_status_types_id
WHERE tst.type = $1;

-- name: GetTasksByPriority :many
SELECT *
FROM tasks
WHERE priority = $1;

-- name: GetTasksByAuthor :many
SELECT *
FROM tasks
WHERE user_id = $1;

-- name: GetTasksByDeadline :many
SELECT *
FROM tasks
WHERE deadline <= $1;

-- name: GetTasksByProject :many
SELECT t.*
FROM tasks AS t
         JOIN project_tasks AS pt ON t.task_id = pt.task_id
WHERE pt.project_id = $1;