ALTER TABLE module_progresses
DROP COLUMN IF EXISTS video_completed_at,
DROP COLUMN IF EXISTS video_watched_percent,
DROP COLUMN IF EXISTS video_watched_seconds,
DROP COLUMN IF EXISTS video_duration_seconds;