-- Add tsvector column for full-text search
ALTER TABLE videos ADD COLUMN search_vector tsvector;

-- Populate existing rows
UPDATE videos SET search_vector =
    setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(description, '')), 'B');

-- Create GIN index for fast full-text search
CREATE INDEX idx_videos_search ON videos USING GIN(search_vector);

-- Create trigger to auto-update search_vector on INSERT or UPDATE
CREATE OR REPLACE FUNCTION videos_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.description, '')), 'B');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_videos_search_vector
    BEFORE INSERT OR UPDATE ON videos
    FOR EACH ROW
    EXECUTE FUNCTION videos_search_vector_update();
