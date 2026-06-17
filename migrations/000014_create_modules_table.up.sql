CREATE TABLE IF NOT EXISTS modules (
    id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,

    module_title VARCHAR(200) NOT NULL,
    module_description TEXT,

    sequence_no INTEGER NOT NULL,
    due_days INTEGER DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_modules_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT uni_modules_course_sequence
        UNIQUE (course_id, sequence_no)
);

CREATE INDEX IF NOT EXISTS idx_modules_course_id
    ON modules(course_id);

CREATE INDEX IF NOT EXISTS idx_modules_is_active
    ON modules(is_active);

CREATE INDEX IF NOT EXISTS idx_modules_sequence_no
    ON modules(sequence_no);