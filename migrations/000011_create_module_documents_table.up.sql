CREATE TABLE IF NOT EXISTS module_documents (
    id BIGSERIAL PRIMARY KEY,

    module_id BIGINT NOT NULL,

    title VARCHAR(200) NOT NULL,
    file_name VARCHAR(255) NOT NULL,

    file_url TEXT NOT NULL,
    public_id VARCHAR(255),

    file_type VARCHAR(20) NOT NULL DEFAULT 'pdf',
    file_size BIGINT NOT NULL,

    page_count INTEGER DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_module_documents_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_module_documents_module_id
    ON module_documents(module_id);

CREATE INDEX IF NOT EXISTS idx_module_documents_file_type
    ON module_documents(file_type);

CREATE INDEX IF NOT EXISTS idx_module_documents_is_active
    ON module_documents(is_active);