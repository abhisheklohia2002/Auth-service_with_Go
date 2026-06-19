ALTER TABLE module_documents
ADD COLUMN thumbnail_name VARCHAR(255),
ADD COLUMN thumbnail_url TEXT,
ADD COLUMN thumbnail_public_id VARCHAR(255),
ADD COLUMN thumbnail_size BIGINT DEFAULT 0;