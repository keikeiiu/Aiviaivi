CREATE TABLE videos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    description TEXT DEFAULT '',
    cover_url VARCHAR(500) NOT NULL DEFAULT '',
    duration INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'uploading'
        CHECK (status IN ('uploading','processing','published','private','deleted')),
    category_id INT REFERENCES categories(id),
    tags TEXT[] DEFAULT '{}',
    view_count BIGINT DEFAULT 0,
    like_count INT DEFAULT 0,
    comment_count INT DEFAULT 0,
    share_count INT DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_videos_user_id ON videos(user_id);
CREATE INDEX idx_videos_status ON videos(status);
CREATE INDEX idx_videos_category ON videos(category_id);
CREATE INDEX idx_videos_published ON videos(published_at DESC);
CREATE INDEX idx_videos_view_count ON videos(view_count DESC);
