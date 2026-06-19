CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,

    course_title VARCHAR(200) NOT NULL,
    course_description TEXT,
    course_type VARCHAR(50) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    created_by_user_id BIGINT NOT NULL,

    CONSTRAINT fk_courses_created_by_user
        FOREIGN KEY (created_by_user_id)
        REFERENCES users(id)
        ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_courses_created_by_user_id
    ON courses(created_by_user_id);

CREATE INDEX IF NOT EXISTS idx_courses_course_type
    ON courses(course_type);

CREATE INDEX IF NOT EXISTS idx_courses_is_active
    ON courses(is_active);