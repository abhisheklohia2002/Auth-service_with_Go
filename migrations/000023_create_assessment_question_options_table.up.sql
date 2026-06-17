CREATE TABLE IF NOT EXISTS assessment_question_options (
    id BIGSERIAL PRIMARY KEY,

    question_id BIGINT NOT NULL,

    option_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_assessment_question_options_question
        FOREIGN KEY (question_id)
        REFERENCES assessment_questions(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_assessment_question_options_question_id
    ON assessment_question_options(question_id);

CREATE INDEX IF NOT EXISTS idx_assessment_question_options_is_correct
    ON assessment_question_options(is_correct);