CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,

    role_name VARCHAR(100) NOT NULL,
    role_type VARCHAR(50) NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uni_roles_role_name UNIQUE (role_name)
);

CREATE INDEX IF NOT EXISTS idx_roles_role_type
    ON roles(role_type);