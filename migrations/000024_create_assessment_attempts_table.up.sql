CREATE TABLE IF NOT EXISTS assessment_attempts (
    id BIGSERIAL PRIMARY KEY,

    assessment_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    attempt_no INTEGER NOT NULL,
    score_obtained INTEGER NOT NULL,
    result_status VARCHAR(50) NOT NULL,

    attempted_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT fk_assessment_attempts_assessment
        FOREIGN KEY (assessment_id)
        REFERENCES assessments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assessment_attempts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT uni_assessment_attempts_user_assessment_attempt_no
        UNIQUE (assessment_id, user_id, attempt_no),

    CONSTRAINT chk_assessment_attempts_attempt_no_positive
        CHECK (attempt_no > 0),

    CONSTRAINT chk_assessment_attempts_score_non_negative
        CHECK (score_obtained >= 0)
);

CREATE INDEX IF NOT EXISTS idx_assessment_attempts_assessment_id
    ON assessment_attempts(assessment_id);

CREATE INDEX IF NOT EXISTS idx_assessment_attempts_user_id
    ON assessment_attempts(user_id);

CREATE INDEX IF NOT EXISTS idx_assessment_attempts_result_status
    ON assessment_attempts(result_status);

CREATE INDEX IF NOT EXISTS idx_assessment_attempts_attempted_at
    ON assessment_attempts(attempted_at);