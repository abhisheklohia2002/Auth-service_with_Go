-- +goose Up
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    role_name VARCHAR(100) NOT NULL UNIQUE,
    role_type VARCHAR(50) NOT NULL,
    description TEXT
);

CREATE INDEX IF NOT EXISTS idx_roles_role_name
ON roles(role_name);

-- +goose Down
DROP INDEX IF EXISTS idx_roles_role_name;

DROP TABLE IF EXISTS roles;