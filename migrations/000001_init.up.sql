CREATE SCHEMA todolist;

CREATE TABLE todolist.users
(
    user_id uuid PRIMARY KEY      DEFAULT uuidv7(),
    version BIGINT       NOT NULL DEFAULT 1,
    email   VARCHAR(255) NOT NULL
        CHECK (email ILIKE '_%@%_._%')
        CHECK (email = BTRIM(email))
);

CREATE TABLE todolist.tasks
(
    task_id      uuid PRIMARY KEY      DEFAULT uuidv7(),
    version      BIGINT       NOT NULL DEFAULT 1,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    completed    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,

    CHECK (
        (completed = FALSE AND completed_at IS NULL)
            OR
        (completed = TRUE AND completed_at IS NOT NULL AND completed_at >= created_at)),

    user_id      uuid         NOT NULL REFERENCES todolist.users (user_id) ON DELETE CASCADE
);

CREATE INDEX idx_tasks_user_id
    ON todolist.tasks (user_id);