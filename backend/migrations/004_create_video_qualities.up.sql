CREATE TABLE video_qualities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    video_id UUID NOT NULL REFERENCES videos(id) ON DELETE CASCADE,
    quality VARCHAR(10) NOT NULL,
    manifest_url VARCHAR(500) NOT NULL,
    bitrate INT,
    file_size BIGINT,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_vq_video_id ON video_qualities(video_id);
