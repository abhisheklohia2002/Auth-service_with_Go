CREATE TABLE IF NOT EXISTS department_training_mappings (
    id BIGSERIAL PRIMARY KEY,

    department_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,

    is_mandatory BOOLEAN NOT NULL DEFAULT TRUE,
    assignment_trigger TEXT,
    active_flag BOOLEAN NOT NULL DEFAULT TRUE,

    created_by_user_id BIGINT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_department_training_mappings_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_department_training_mappings_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_department_training_mappings_created_by
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT uni_department_training_mappings_department_course
        UNIQUE (department_id, course_id)
);

CREATE INDEX IF NOT EXISTS idx_department_training_mappings_department_id
    ON department_training_mappings(department_id);

CREATE INDEX IF NOT EXISTS idx_department_training_mappings_course_id
    ON department_training_mappings(course_id);

CREATE INDEX IF NOT EXISTS idx_department_training_mappings_created_by_user_id
    ON department_training_mappings(created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_department_training_mappings_active_flag
    ON department_training_mappings(active_flag);

CREATE INDEX IF NOT EXISTS idx_department_training_mappings_is_mandatory
    ON department_training_mappings(is_mandatory);