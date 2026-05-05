package model

import (
	"context"
	"database/sql"
	"time"
)

type CreatorOverview struct {
	TotalViews       int64 `json:"total_views"`
	TotalLikes       int64 `json:"total_likes"`
	TotalComments    int64 `json:"total_comments"`
	TotalVideos      int64 `json:"total_videos"`
	TotalSubscribers int64 `json:"total_subscribers"`
	TotalWatchTime   int64 `json:"total_watch_time_seconds"`
}

type VideoStats struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	CoverURL     string    `json:"cover_url"`
	ViewCount    int64     `json:"view_count"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"`
	ShareCount   int       `json:"share_count"`
	PublishedAt  *time.Time `json:"published_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type VideoStatsDetail struct {
	VideoStats
	DailyViews []DailyMetric `json:"daily_views"`
}

type DailyMetric struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

func GetCreatorOverview(ctx context.Context, db *sql.DB, userID string) (CreatorOverview, error) {
	var o CreatorOverview

	err := db.QueryRowContext(ctx, `
	SELECT
	    COALESCE(SUM(v.view_count), 0),
	    COALESCE(SUM(v.like_count), 0),
	    COALESCE(SUM(v.comment_count), 0),
	    COALESCE(COUNT(v.id), 0),
	    COALESCE((SELECT count(*) FROM follows WHERE creator_id = $1), 0),
	    COALESCE(SUM(v.duration * v.view_count), 0)
	FROM videos v
	WHERE v.user_id = $1 AND v.status != 'deleted'
	`, userID).Scan(
		&o.TotalViews, &o.TotalLikes, &o.TotalComments,
		&o.TotalVideos, &o.TotalSubscribers, &o.TotalWatchTime,
	)
	if err != nil {
		return CreatorOverview{}, err
	}
	return o, nil
}

func GetCreatorVideoStats(ctx context.Context, db *sql.DB, userID string) ([]VideoStats, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT id, title, cover_url, view_count, like_count, comment_count, share_count, published_at, created_at
	FROM videos
	WHERE user_id = $1 AND status != 'deleted'
	ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []VideoStats
	for rows.Next() {
		var s VideoStats
		if err := rows.Scan(&s.ID, &s.Title, &s.CoverURL, &s.ViewCount, &s.LikeCount, &s.CommentCount, &s.ShareCount, &s.PublishedAt, &s.CreatedAt); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}
