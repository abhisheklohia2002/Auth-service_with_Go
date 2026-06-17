CREATE TABLE IF NOT EXISTS notification_recipients (
    id BIGSERIAL PRIMARY KEY,

    notification_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,

    read_status BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ NULL,
    delivered_at TIMESTAMPTZ NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_notification_recipients_notification
        FOREIGN KEY (notification_id)
        REFERENCES notifications(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_notification_recipients_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT idx_notification_user
        UNIQUE (notification_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_notification_recipients_notification_id
    ON notification_recipients(notification_id);

CREATE INDEX IF NOT EXISTS idx_notification_recipients_user_id
    ON notification_recipients(user_id);

CREATE INDEX IF NOT EXISTS idx_notification_recipients_read_status
    ON notification_recipients(read_status);

CREATE INDEX IF NOT EXISTS idx_notification_recipients_read_at
    ON notification_recipients(read_at);

CREATE INDEX IF NOT EXISTS idx_notification_recipients_delivered_at
    ON notification_recipients(delivered_at);