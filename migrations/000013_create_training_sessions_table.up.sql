CREATE TABLE IF NOT EXISTS training_sessions (
    id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,
    module_id BIGINT NULL,
    created_by_user_id BIGINT NOT NULL,

    session_title VARCHAR(255) NOT NULL,
    session_type VARCHAR(50) NOT NULL,

    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,

    location VARCHAR(255),
    meeting_link VARCHAR(500),

    is_mandatory BOOLEAN NOT NULL DEFAULT TRUE,
    status VARCHAR(50) DEFAULT 'scheduled',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_training_sessions_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_training_sessions_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_training_sessions_created_by_user
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_training_sessions_time_range
        CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_training_sessions_course_id
    ON training_sessions(course_id);

CREATE INDEX IF NOT EXISTS idx_training_sessions_module_id
    ON training_sessions(module_id);

CREATE INDEX IF NOT EXISTS idx_training_sessions_created_by_user_id
    ON training_sessions(created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_training_sessions_session_type
    ON training_sessions(session_type);

CREATE INDEX IF NOT EXISTS idx_training_sessions_status
    ON training_sessions(status);

CREATE INDEX IF NOT EXISTS idx_training_sessions_start_time
    ON training_sessions(start_time);

CREATE INDEX IF NOT EXISTS idx_training_sessions_end_time
    ON training_sessions(end_time);

CREATE INDEX IF NOT EXISTS idx_training_sessions_is_mandatory
    ON training_sessions(is_mandatory);