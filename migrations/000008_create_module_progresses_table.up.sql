CREATE TABLE IF NOT EXISTS module_progresses (
    id BIGSERIAL PRIMARY KEY,

    assignment_id BIGINT NOT NULL,
    module_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    status VARCHAR(30) DEFAULT 'pending',

    started_at TIMESTAMPTZ NULL,
    completed_at TIMESTAMPTZ NULL,

    CONSTRAINT fk_module_progresses_assignment
        FOREIGN KEY (assignment_id)
        REFERENCES training_assignments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_module_progresses_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_module_progresses_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT idx_assignment_module
        UNIQUE (assignment_id, module_id)
);

CREATE INDEX IF NOT EXISTS idx_module_progresses_assignment_id
    ON module_progresses(assignment_id);

CREATE INDEX IF NOT EXISTS idx_module_progresses_module_id
    ON module_progresses(module_id);

CREATE INDEX IF NOT EXISTS idx_module_progresses_user_id
    ON module_progresses(user_id);

CREATE INDEX IF NOT EXISTS idx_module_progresses_status
    ON module_progresses(status);