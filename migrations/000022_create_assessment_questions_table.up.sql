CREATE TABLE IF NOT EXISTS assessment_questions (
    id BIGSERIAL PRIMARY KEY,

    assessment_id BIGINT NOT NULL,

    question_text TEXT NOT NULL,
    question_type VARCHAR(50) NOT NULL,

    marks INTEGER NOT NULL DEFAULT 1,
    sequence_no INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_assessment_questions_assessment
        FOREIGN KEY (assessment_id)
        REFERENCES assessments(id)
        ON DELETE CASCADE,

    CONSTRAINT uni_assessment_questions_assessment_sequence
        UNIQUE (assessment_id, sequence_no),

    CONSTRAINT chk_assessment_questions_marks_positive
        CHECK (marks > 0)
);

CREATE INDEX IF NOT EXISTS idx_assessment_questions_assessment_id
    ON assessment_questions(assessment_id);

CREATE INDEX IF NOT EXISTS idx_assessment_questions_question_type
    ON assessment_questions(question_type);

CREATE INDEX IF NOT EXISTS idx_assessment_questions_is_active
    ON assessment_questions(is_active);

CREATE INDEX IF NOT EXISTS idx_assessment_questions_sequence_no
    ON assessment_questions(sequence_no);