CREATE TABLE module_videos (
    video_id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,
    module_id BIGINT NOT NULL,

    title VARCHAR(255),
    video_name VARCHAR(255),
    video_url TEXT NOT NULL,
    video_public_id VARCHAR(255) NOT NULL,
    video_size BIGINT DEFAULT 0,
    video_type VARCHAR(50) DEFAULT 'video',

    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX idx_module_videos_module_id
ON module_videos(module_id);

CREATE INDEX idx_module_videos_course_id
ON module_videos(course_id);

CREATE INDEX idx_module_videos_is_active
ON module_videos(is_active);