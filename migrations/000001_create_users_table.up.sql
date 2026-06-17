CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,

    full_name VARCHAR(150) NOT NULL,
    email VARCHAR(150) NOT NULL,
    password TEXT NOT NULL,

    role_id BIGINT NOT NULL,

    status VARCHAR(30) DEFAULT 'active',

    manager_id BIGINT NULL,

    joining_date TIMESTAMPTZ NOT NULL,
    employee_code VARCHAR(50) NOT NULL,

    department_id BIGINT NULL,

    CONSTRAINT uni_users_email UNIQUE (email),
    CONSTRAINT uni_users_employee_code UNIQUE (employee_code),

    CONSTRAINT fk_users_role
        FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_users_manager
        FOREIGN KEY (manager_id)
        REFERENCES users(id)
        ON DELETE SET NULL,

    CONSTRAINT fk_users_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_users_role_id
    ON users(role_id);

CREATE INDEX IF NOT EXISTS idx_users_manager_id
    ON users(manager_id);

CREATE INDEX IF NOT EXISTS idx_users_department_id
    ON users(department_id);

CREATE INDEX IF NOT EXISTS idx_users_status
    ON users(status);