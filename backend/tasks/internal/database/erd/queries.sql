-- name: CreateTask :one
INSERT INTO tasks (title, description, priority, author_name, created_at, deadline)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: CreateNewTaskStatus :one
INSERT INTO task_status_types (type)
VALUES ($1)
RETURNING *;

-- name: SetTaskDeadline :exec
UPDATE tasks
SET deadline = $2
WHERE task_id = $1;

-- name: SetTaskStatus :exec
UPDATE task_statuses
SET task_status_id = (SELECT task_status_id FROM task_status_types WHERE $2 = type)
WHERE task_id = $1;

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
    author_name = $5,
    created_at  = $6,
    deadline    = $7
WHERE task_id = $1;

-- name: GetTaskStatusByTaskId :one
SELECT t.type
FROM task_statuses AS ts
         JOIN task_status_types AS t ON ts.task_status_id = t.task_status_id
WHERE ts.task_id = $1;

-- name: GetAllTaskStatuses :many
SELECT *
FROM task_status_types;

-- name: GetTasksByStatus :many
SELECT t.*
FROM tasks AS t
         JOIN task_statuses AS ts ON t.task_id = ts.task_id
         JOIN task_status_types AS tst ON ts.task_status_id = tst.task_status_id
WHERE tst.type = $1;

-- name: GetTasksByPriority :many
SELECT * FROM tasks
WHERE priority = $1;

-- name: GetTasksByAuthor :many
SELECT * FROM tasks
WHERE author_name = $1;

-- name: GetTasksByDeadline :many
SELECT * FROM tasks
WHERE deadline <= $1;