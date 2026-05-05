CREATE TABLE watch_history (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    progress FLOAT NOT NULL DEFAULT 0,     -- seconds watched
    duration FLOAT NOT NULL DEFAULT 0,     -- video duration at time of recording
    watched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(user_id, video_id)              -- one entry per user per video
);

CREATE INDEX idx_watch_history_user ON watch_history(user_id, watched_at DESC);
CREATE INDEX idx_watch_history_video ON watch_history(video_id);
