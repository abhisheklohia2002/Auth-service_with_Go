CREATE TABLE IF NOT EXISTS training_assignments (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    assigned_by_user_id BIGINT NOT NULL,

    assignment_source VARCHAR(100),
    is_mandatory BOOLEAN NOT NULL DEFAULT FALSE,

    assigned_date TIMESTAMPTZ NOT NULL,
    due_date TIMESTAMPTZ NULL,
    completion_date TIMESTAMPTZ NULL,

    status VARCHAR(50) DEFAULT 'assigned',
    improvement_status VARCHAR(50),

    department_id BIGINT NULL,

    CONSTRAINT fk_training_assignments_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_training_assignments_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_training_assignments_assigned_by_user
        FOREIGN KEY (assigned_by_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_training_assignments_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_training_assignments_user_id
    ON training_assignments(user_id);

CREATE INDEX IF NOT EXISTS idx_training_assignments_course_id
    ON training_assignments(course_id);

CREATE INDEX IF NOT EXISTS idx_training_assignments_assigned_by_user_id
    ON training_assignments(assigned_by_user_id);

CREATE INDEX IF NOT EXISTS idx_training_assignments_department_id
    ON training_assignments(department_id);

CREATE INDEX IF NOT EXISTS idx_training_assignments_status
    ON training_assignments(status);

CREATE INDEX IF NOT EXISTS idx_training_assignments_due_date
    ON training_assignments(due_date);

CREATE INDEX IF NOT EXISTS idx_training_assignments_is_mandatory
    ON training_assignments(is_mandatory);