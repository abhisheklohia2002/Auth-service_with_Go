CREATE TABLE IF NOT EXISTS departments (
    id BIGSERIAL PRIMARY KEY,

    entity_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_departments_entity
        FOREIGN KEY (entity_id)
        REFERENCES entities(id)
        ON DELETE CASCADE,

    CONSTRAINT idx_entity_department_name
        UNIQUE (entity_id, name)
);

CREATE INDEX IF NOT EXISTS idx_departments_entity_id
    ON departments(entity_id);

CREATE INDEX IF NOT EXISTS idx_departments_is_active
    ON departments(is_active);