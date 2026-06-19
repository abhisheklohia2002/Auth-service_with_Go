CREATE TABLE IF NOT EXISTS training_mappings (
    id BIGSERIAL PRIMARY KEY,

    role_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,

    is_mandatory BOOLEAN NOT NULL DEFAULT FALSE,
    assignment_trigger VARCHAR(100),
    active_flag BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_training_mappings_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_training_mappings_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_training_mappings_role_id
    ON training_mappings(role_id);

CREATE INDEX IF NOT EXISTS idx_training_mappings_course_id
    ON training_mappings(course_id);

CREATE INDEX IF NOT EXISTS idx_training_mappings_active_flag
    ON training_mappings(active_flag);

CREATE INDEX IF NOT EXISTS idx_training_mappings_is_mandatory
    ON training_mappings(is_mandatory);