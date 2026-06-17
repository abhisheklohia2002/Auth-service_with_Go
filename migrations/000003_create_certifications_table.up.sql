CREATE TABLE IF NOT EXISTS certifications (
    id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,
    certification_name VARCHAR(200) NOT NULL,
    validity_days INTEGER NOT NULL,
    rule_id BIGINT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_certifications_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_certifications_rule
        FOREIGN KEY (rule_id)
        REFERENCES certification_rules(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_certifications_course_id
    ON certifications(course_id);

CREATE INDEX IF NOT EXISTS idx_certifications_rule_id
    ON certifications(rule_id);

CREATE INDEX IF NOT EXISTS idx_certifications_is_active
    ON certifications(is_active);