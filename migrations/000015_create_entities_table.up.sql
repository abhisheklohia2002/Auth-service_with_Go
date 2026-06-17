CREATE TABLE IF NOT EXISTS entities (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    entity_type TEXT,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uni_entities_name UNIQUE (name)
);

CREATE INDEX IF NOT EXISTS idx_entities_entity_type
    ON entities(entity_type);

CREATE INDEX IF NOT EXISTS idx_entities_is_active
    ON entities(is_active);