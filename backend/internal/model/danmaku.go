package model

import (
	"context"
	"database/sql"
	"time"
)

type Danmaku struct {
	ID        int64     `json:"id"`
	VideoID   string    `json:"video_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	VideoTime float64   `json:"video_time"`
	Color     string    `json:"color"`
	FontSize  string    `json:"font_size"`
	Mode      string    `json:"mode"`
	CreatedAt time.Time `json:"created_at"`
	// Joined
	Username  string `json:"username"`
}

func CreateDanmaku(ctx context.Context, db *sql.DB, videoID string, userID string, content string, videoTime float64, color string, fontSize string, mode string) (Danmaku, error) {
	var d Danmaku
	err := db.QueryRowContext(ctx, `
	INSERT INTO danmaku(video_id, user_id, content, video_time, color, font_size, mode)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id, video_id, user_id, content, video_time, color, font_size, mode, created_at
	`, videoID, userID, content, videoTime, color, fontSize, mode).Scan(
		&d.ID, &d.VideoID, &d.UserID, &d.Content, &d.VideoTime,
		&d.Color, &d.FontSize, &d.Mode, &d.CreatedAt,
	)
	return d, err
}

func GetDanmakuByTimeRange(ctx context.Context, db *sql.DB, videoID string, tStart float64, tEnd float64) ([]Danmaku, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT d.id, d.video_id, d.user_id, d.content, d.video_time, d.color, d.font_size, d.mode, d.created_at, u.username
	FROM danmaku d JOIN users u ON d.user_id = u.id
	WHERE d.video_id = $1 AND d.video_time >= $2 AND d.video_time <= $3
	ORDER BY d.video_time ASC, d.created_at ASC
	LIMIT 500
	`, videoID, tStart, tEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Danmaku
	for rows.Next() {
		var d Danmaku
		if err := rows.Scan(
			&d.ID, &d.VideoID, &d.UserID, &d.Content, &d.VideoTime,
			&d.Color, &d.FontSize, &d.Mode, &d.CreatedAt, &d.Username,
		); err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}
