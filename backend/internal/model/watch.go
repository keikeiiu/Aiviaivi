package model

import (
	"context"
	"database/sql"
	"time"
)

type WatchEntry struct {
	ID        int64     `json:"id"`
	UserID    string    `json:"user_id"`
	VideoID   string    `json:"video_id"`
	Progress  float64   `json:"progress"`
	Duration  float64   `json:"duration"`
	WatchedAt time.Time `json:"watched_at"`
	// Joined
	Title     string `json:"title"`
	CoverURL  string `json:"cover_url"`
}

func RecordWatch(ctx context.Context, db *sql.DB, userID string, videoID string, progress float64, duration float64) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO watch_history(user_id, video_id, progress, duration, watched_at)
	VALUES ($1, $2, $3, $4, now())
	ON CONFLICT (user_id, video_id) DO UPDATE
	SET progress = EXCLUDED.progress,
	    duration = EXCLUDED.duration,
	    watched_at = now()
	`, userID, videoID, progress, duration)
	return err
}

func GetWatchHistory(ctx context.Context, db *sql.DB, userID string, page int, size int) ([]WatchEntry, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	offset := (page - 1) * size

	var total int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM watch_history WHERE user_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := db.QueryContext(ctx, `
	SELECT w.id, w.user_id, w.video_id, w.progress, w.duration, w.watched_at,
	       v.title, v.cover_url
	FROM watch_history w
	JOIN videos v ON w.video_id = v.id
	WHERE w.user_id = $1
	ORDER BY w.watched_at DESC
	LIMIT $2 OFFSET $3
	`, userID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var entries []WatchEntry
	for rows.Next() {
		var e WatchEntry
		if err := rows.Scan(
			&e.ID, &e.UserID, &e.VideoID, &e.Progress, &e.Duration, &e.WatchedAt,
			&e.Title, &e.CoverURL,
		); err != nil {
			return nil, 0, err
		}
		entries = append(entries, e)
	}
	return entries, total, rows.Err()
}
