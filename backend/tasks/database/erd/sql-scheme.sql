CREATE TYPE task_priority as ENUM ('low', 'medium', 'high', 'really high');

CREATE TABLE tasks
(
    task_id     BIGINT        NOT NULL GENERATED ALWAYS AS IDENTITY,
    title       TEXT          NOT NULL,
    description TEXT,
    priority    task_priority NOT NULL,
    author_name TEXT          NOT NULL,
    created_at   TIMESTAMP,
    deadline    TIMESTAMP
);

CREATE TABLE task_status_types
(
    task_status_id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY,
    type           TEXT   NOT NULL UNIQUE
);
CREATE UNIQUE INDEX idx_task_status_type ON task_status_types (type);

CREATE TABLE task_statuses
(
    task_id        BIGINT NOT NULL,
    task_status_id BIGINT NOT NULL,

    CONSTRAINT tasks_id_fk FOREIGN KEY (task_id) REFERENCES tasks (task_id),
    CONSTRAINT task_status_id_fk FOREIGN KEY (task_status_id) REFERENCES task_status_types (task_status_id)
);