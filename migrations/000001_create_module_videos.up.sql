CREATE TABLE IF NOT EXISTS module_videos (
    video_id BIGSERIAL PRIMARY KEY,

    course_id BIGINT NOT NULL,
    module_id BIGINT NOT NULL,

    title TEXT,
    video_name TEXT,
    video_url TEXT NOT NULL,
    video_public_id TEXT NOT NULL,
    video_size BIGINT DEFAULT 0,
    video_type TEXT,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_module_videos_course_id
ON module_videos(course_id);

CREATE INDEX IF NOT EXISTS idx_module_videos_module_id
ON module_videos(module_id);

CREATE INDEX IF NOT EXISTS idx_module_videos_is_active
ON module_videos(is_active);