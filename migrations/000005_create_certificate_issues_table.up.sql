CREATE TABLE IF NOT EXISTS certificate_issues (
    id BIGSERIAL PRIMARY KEY,

    certification_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    training_assignment_id BIGINT NULL,

    issued_on TIMESTAMPTZ NOT NULL,
    expiry_date TIMESTAMPTZ NOT NULL,

    certificate_number VARCHAR(100) NOT NULL,
    issue_status VARCHAR(50) DEFAULT 'issued',
    pdf_path VARCHAR(500),

    CONSTRAINT uni_certificate_issues_certificate_number
        UNIQUE (certificate_number),

    CONSTRAINT fk_certificate_issues_certification
        FOREIGN KEY (certification_id)
        REFERENCES certifications(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_certificate_issues_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_certificate_issues_training_assignment
        FOREIGN KEY (training_assignment_id)
        REFERENCES training_assignments(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_certificate_issues_certification_id
    ON certificate_issues(certification_id);

CREATE INDEX IF NOT EXISTS idx_certificate_issues_user_id
    ON certificate_issues(user_id);

CREATE INDEX IF NOT EXISTS idx_certificate_issues_training_assignment_id
    ON certificate_issues(training_assignment_id);

CREATE INDEX IF NOT EXISTS idx_certificate_issues_issue_status
    ON certificate_issues(issue_status);

CREATE INDEX IF NOT EXISTS idx_certificate_issues_expiry_date
    ON certificate_issues(expiry_date);