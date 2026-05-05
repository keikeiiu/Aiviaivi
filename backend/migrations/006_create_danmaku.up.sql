CREATE TABLE danmaku (
    id BIGSERIAL PRIMARY KEY,
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    video_time FLOAT NOT NULL,
    color VARCHAR(7) DEFAULT '#FFFFFF',
    font_size VARCHAR(10) DEFAULT 'medium'
        CHECK (font_size IN ('small','medium','large')),
    mode VARCHAR(10) DEFAULT 'scroll'
        CHECK (mode IN ('scroll','top','bottom')),
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_danmaku_video_time ON danmaku(video_id, video_time);
