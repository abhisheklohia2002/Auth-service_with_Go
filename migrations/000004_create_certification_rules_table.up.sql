CREATE TABLE IF NOT EXISTS certification_rules (
    id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,

    issue_on_course_completion BOOLEAN NOT NULL DEFAULT TRUE,
    minimum_score_required INTEGER NOT NULL,
    validity_days INTEGER NOT NULL,
    renewal_required BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_certification_rules_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_certification_rules_course_id
    ON certification_rules(course_id);