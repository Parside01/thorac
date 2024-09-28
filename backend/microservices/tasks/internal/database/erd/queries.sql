-- name: CreateProject :one
INSERT INTO projects (title, description, user_id, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;

-- name: AddUserToProject :one
INSERT INTO users_projects (user_id, project_id)
VALUES ($1, $2)
RETURNING *;

-- name: CreateTask :one
WITH new_task AS (
    INSERT INTO tasks (title, description, priority, user_id, created_at, deadline)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING *),
new_project_task AS (
    INSERT INTO project_tasks (task_id, project_id)
    SELECT task_id, $7 FROM new_task
    RETURNING *
)
INSERT
INTO users_tasks (user_id, task_id)
SELECT $8, task_id
FROM new_task
RETURNING *;

-- name: AddTaskToProject :exec
WITH new_project_task AS (
    INSERT INTO project_tasks (task_id, project_id)
        VALUES ($1, $2)
        RETURNING *
)
INSERT INTO users_tasks (user_id, task_id)
SELECT $3, task_id FROM new_project_task;

-- name: CreateNewTaskStatus :one
INSERT INTO task_status_types (owner_project_id, type)
VALUES ($1, $2)
RETURNING *;

-- name: SetTaskDeadline :exec
UPDATE tasks
SET deadline = $2
WHERE task_id = $1;

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

-- name: GetTaskById :one
SELECT *
FROM tasks
WHERE task_id = $1;

-- name: DeleteTask :exec
DELETE
FROM tasks
WHERE task_id = $1;

-- name: UpdateTask :exec
UPDATE tasks
SET title       = $2,
    description = $3,
    priority    = $4,
    user_id = $5,
    created_at  = $6,
    deadline    = $7
WHERE task_id = $1;

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