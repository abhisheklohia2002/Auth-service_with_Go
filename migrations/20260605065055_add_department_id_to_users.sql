-- +goose Up
ALTER TABLE users
ADD COLUMN IF NOT EXISTS department_id BIGINT;

CREATE INDEX IF NOT EXISTS idx_users_department_id
ON users(department_id);

ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_department;

ALTER TABLE users
ADD CONSTRAINT fk_users_department
FOREIGN KEY (department_id)
REFERENCES departments(id)
ON DELETE SET NULL;

-- +goose Down
ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_department;

DROP INDEX IF EXISTS idx_users_department_id;

ALTER TABLE users
DROP COLUMN IF EXISTS department_id;