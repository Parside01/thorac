-- Define a new enum type for task priorities
CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high', 'really high');

-- Create the projects table
CREATE TABLE projects
(
    project_id  BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY,
    title       TEXT   NOT NULL,
    description TEXT,
    author_name TEXT   NOT NULL,
    created_at  TIMESTAMP,
    PRIMARY KEY (project_id)
);

-- Create the tasks table
CREATE TABLE tasks
(
    task_id     BIGINT        NOT NULL GENERATED ALWAYS AS IDENTITY,
    title       TEXT          NOT NULL,
    description TEXT,
    priority    task_priority NOT NULL,
    author_name TEXT          NOT NULL,
    created_at  TIMESTAMP,
    deadline    TIMESTAMP,
    PRIMARY KEY (task_id)
);

-- Create the project_tasks table
CREATE TABLE project_tasks
(
    task_id    BIGINT NOT NULL,
    project_id BIGINT NOT NULL,
    CONSTRAINT project_tasks_task_fk FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    CONSTRAINT project_tasks_project_fk FOREIGN KEY (project_id) REFERENCES projects (project_id)
);

-- Create the task_status_types table
CREATE TABLE task_status_types
(
    task_status_types_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY,
    owner_project_id     BIGINT NOT NULL,
    type                 TEXT   NOT NULL,
    PRIMARY KEY (task_status_types_id),
    UNIQUE (owner_project_id, type),
    CONSTRAINT task_status_types_project_fk FOREIGN KEY (owner_project_id) REFERENCES projects (project_id)
);

-- Create the task_statuses table
CREATE TABLE task_statuses
(
    task_id        BIGINT NOT NULL,
    task_status_id BIGINT NOT NULL,
    PRIMARY KEY (task_id, task_status_id),
    CONSTRAINT task_statuses_task_fk FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    CONSTRAINT task_statuses_status_fk FOREIGN KEY (task_status_id) REFERENCES task_status_types (task_status_types_id)
);