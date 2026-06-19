CREATE TABLE IF NOT EXISTS notifications (
    id BIGSERIAL PRIMARY KEY,

    created_by_id BIGINT NOT NULL,

    assignment_id BIGINT NULL,
    certificate_issue_id BIGINT NULL,

    notification_type VARCHAR(100) NOT NULL,
    title VARCHAR(255),
    message TEXT NOT NULL,

    target_audience VARCHAR(100),
    sent_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notifications_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_notifications_assignment
        FOREIGN KEY (assignment_id)
        REFERENCES training_assignments(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_notifications_certificate_issue
        FOREIGN KEY (certificate_issue_id)
        REFERENCES certificate_issues(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_notifications_created_by_id
    ON notifications(created_by_id);

CREATE INDEX IF NOT EXISTS idx_notifications_assignment_id
    ON notifications(assignment_id);

CREATE INDEX IF NOT EXISTS idx_notifications_certificate_issue_id
    ON notifications(certificate_issue_id);

CREATE INDEX IF NOT EXISTS idx_notifications_notification_type
    ON notifications(notification_type);

CREATE INDEX IF NOT EXISTS idx_notifications_target_audience
    ON notifications(target_audience);

CREATE INDEX IF NOT EXISTS idx_notifications_sent_at
    ON notifications(sent_at);