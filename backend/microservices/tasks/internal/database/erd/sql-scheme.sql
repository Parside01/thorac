CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'really high');

CREATE TABLE IF NOT EXISTS projects
(
    project_id  TEXT NOT NULL GENERATED ALWAYS AS IDENTITY,
    title       TEXT   NOT NULL,
    description TEXT,
    user_id     TEXT   NOT NULL, --- Автор проекта.
    created_at  TIMESTAMP,
    PRIMARY KEY (project_id)
);

CREATE TABLE IF NOT EXISTS tasks
(
    task_id     TEXT        NOT NULL GENERATED ALWAYS AS IDENTITY,
    title       TEXT          NOT NULL,
    description TEXT,
    priority    task_priority NOT NULL,
    user_id     TEXT          NOT NULL, --- Автор задачи.
    created_at  TIMESTAMP,
    deadline    TIMESTAMP,
    PRIMARY KEY (task_id)
);

CREATE TABLE IF NOT EXISTS project_tasks --- Все задачи, №задачи - проект в котором она находится.
(
    task_id    TEXT NOT NULL,
    project_id TEXT NOT NULL,
    CONSTRAINT project_tasks_task_fk FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    CONSTRAINT project_tasks_project_fk FOREIGN KEY (project_id) REFERENCES projects (project_id)
);

CREATE TABLE IF NOT EXISTS task_status_types
(
    task_status_types_id TEXT NOT NULL GENERATED ALWAYS AS IDENTITY,
    owner_project_id     TEXT NOT NULL,
    type                 TEXT   NOT NULL,
    PRIMARY KEY (task_status_types_id),
    UNIQUE (type),
    CONSTRAINT task_status_types_project_fk FOREIGN KEY (owner_project_id) REFERENCES projects (project_id)
);

CREATE TABLE IF NOT EXISTS task_statuses
(
    task_id        TEXT NOT NULL,
    task_status_id TEXT NOT NULL,
    PRIMARY KEY (task_id, task_status_id),
    CONSTRAINT task_statuses_task_fk FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    CONSTRAINT task_statuses_status_fk FOREIGN KEY (task_status_id) REFERENCES task_status_types (task_status_types_id)
);

CREATE TABLE IF NOT EXISTS users_projects --- Все юзеры и в каких проектах они находятся.
(
    user_id    TEXT NOT NULL,
    project_id TEXT NOT NULL,

    CONSTRAINT project_tasks_project_fk FOREIGN KEY (project_id) REFERENCES projects (project_id)
);

CREATE TABLE IF NOT EXISTS users_tasks --- Все юзеры и какие таски у них есть.
(
    user_id TEXT NOT NULL,
    task_id TEXT NOT NULL,

    CONSTRAINT task_statuses_task_fk FOREIGN KEY (task_id) REFERENCES tasks (task_id)
)