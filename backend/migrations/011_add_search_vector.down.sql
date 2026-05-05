DROP TRIGGER IF EXISTS trg_videos_search_vector ON videos;
DROP FUNCTION IF EXISTS videos_search_vector_update();
DROP INDEX IF EXISTS idx_videos_search;
ALTER TABLE videos DROP COLUMN IF EXISTS search_vector;
