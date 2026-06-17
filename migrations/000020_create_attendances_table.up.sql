CREATE TABLE IF NOT EXISTS attendances (
    id BIGSERIAL PRIMARY KEY,

    session_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    marked_by_user_id BIGINT NULL,

    status VARCHAR(50) NOT NULL,

    check_in_time TIMESTAMPTZ NULL,
    check_out_time TIMESTAMPTZ NULL,

    duration_minutes INTEGER NOT NULL DEFAULT 0,
    attendance_source VARCHAR(50) DEFAULT 'manual',

    remarks VARCHAR(500),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_attendances_session
        FOREIGN KEY (session_id)
        REFERENCES training_sessions(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_attendances_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_attendances_marked_by_user
        FOREIGN KEY (marked_by_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT idx_session_user
        UNIQUE (session_id, user_id),

    CONSTRAINT chk_attendances_duration_non_negative
        CHECK (duration_minutes >= 0),

    CONSTRAINT chk_attendances_checkout_after_checkin
        CHECK (
            check_out_time IS NULL
            OR check_in_time IS NULL
            OR check_out_time >= check_in_time
        )
);

CREATE INDEX IF NOT EXISTS idx_attendances_session_id
    ON attendances(session_id);

CREATE INDEX IF NOT EXISTS idx_attendances_user_id
    ON attendances(user_id);

CREATE INDEX IF NOT EXISTS idx_attendances_marked_by_user_id
    ON attendances(marked_by_user_id);

CREATE INDEX IF NOT EXISTS idx_attendances_status
    ON attendances(status);

CREATE INDEX IF NOT EXISTS idx_attendances_attendance_source
    ON attendances(attendance_source);

CREATE INDEX IF NOT EXISTS idx_attendances_check_in_time
    ON attendances(check_in_time);