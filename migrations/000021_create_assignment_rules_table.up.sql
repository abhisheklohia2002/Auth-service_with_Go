CREATE TABLE IF NOT EXISTS assignment_rules (
    assignment_rule_id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,

    rule_name VARCHAR(255) NOT NULL,
    trigger_event VARCHAR(255) NOT NULL,

    due_days INTEGER NOT NULL DEFAULT 14,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_assignment_rules_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_assignment_rules_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE,

    CONSTRAINT uni_assignment_rules_course_role_trigger
        UNIQUE (course_id, role_id, trigger_event),

    CONSTRAINT chk_assignment_rules_due_days_non_negative
        CHECK (due_days >= 0)
);

CREATE INDEX IF NOT EXISTS idx_assignment_rules_course_id
    ON assignment_rules(course_id);

CREATE INDEX IF NOT EXISTS idx_assignment_rules_role_id
    ON assignment_rules(role_id);

CREATE INDEX IF NOT EXISTS idx_assignment_rules_trigger_event
    ON assignment_rules(trigger_event);

CREATE INDEX IF NOT EXISTS idx_assignment_rules_is_active
    ON assignment_rules(is_active);