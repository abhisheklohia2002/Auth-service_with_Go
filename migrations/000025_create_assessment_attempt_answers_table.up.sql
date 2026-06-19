CREATE TABLE IF NOT EXISTS assessment_attempt_answers (
    id BIGSERIAL PRIMARY KEY,

    attempt_id BIGINT NOT NULL,
    question_id BIGINT NOT NULL,

    selected_option_ids TEXT,
    text_answer TEXT,

    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    marks_awarded INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT fk_assessment_attempt_answers_attempt
        FOREIGN KEY (attempt_id)
        REFERENCES assessment_attempts(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assessment_attempt_answers_question
        FOREIGN KEY (question_id)
        REFERENCES assessment_questions(id)
        ON DELETE CASCADE,

    CONSTRAINT uni_assessment_attempt_answers_attempt_question
        UNIQUE (attempt_id, question_id),

    CONSTRAINT chk_assessment_attempt_answers_marks_non_negative
        CHECK (marks_awarded >= 0)
);

CREATE INDEX IF NOT EXISTS idx_assessment_attempt_answers_attempt_id
    ON assessment_attempt_answers(attempt_id);

CREATE INDEX IF NOT EXISTS idx_assessment_attempt_answers_question_id
    ON assessment_attempt_answers(question_id);

CREATE INDEX IF NOT EXISTS idx_assessment_attempt_answers_is_correct
    ON assessment_attempt_answers(is_correct);