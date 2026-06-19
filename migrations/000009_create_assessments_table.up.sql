CREATE TABLE IF NOT EXISTS assessments (
    id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,
    module_id BIGINT NULL,

    assessment_title VARCHAR(200) NOT NULL,
    assessment_type VARCHAR(100) NOT NULL,

    max_score INTEGER NOT NULL,
    passing_score INTEGER NOT NULL,

    rule_id BIGINT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_assessments_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assessments_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_assessments_rule
        FOREIGN KEY (rule_id)
        REFERENCES assessment_rules(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_assessments_course_id
    ON assessments(course_id);

CREATE INDEX IF NOT EXISTS idx_assessments_module_id
    ON assessments(module_id);

CREATE INDEX IF NOT EXISTS idx_assessments_rule_id
    ON assessments(rule_id);

CREATE INDEX IF NOT EXISTS idx_assessments_assessment_type
    ON assessments(assessment_type);

CREATE INDEX IF NOT EXISTS idx_assessments_is_active
    ON assessments(is_active);