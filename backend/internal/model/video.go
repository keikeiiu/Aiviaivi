package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type Video struct {
	ID           string
	UserID       string
	Title        string
	Description  string
	CoverURL     string
	Duration     int
	Status       string
	CategoryID   *int
	Tags         []string
	ViewCount    int64
	LikeCount    int
	CommentCount int
	ShareCount   int
	PublishedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// Joined fields
	Username  string `json:"-"`
	AvatarURL string `json:"-"`
}

type VideoDetail struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	CoverURL     string          `json:"cover_url"`
	Duration     int             `json:"duration"`
	Status       string          `json:"status"`
	CategoryID   *int            `json:"category_id"`
	Tags         []string        `json:"tags"`
	ViewCount    int64           `json:"view_count"`
	LikeCount    int             `json:"like_count"`
	CommentCount int             `json:"comment_count"`
	ShareCount   int             `json:"share_count"`
	PublishedAt  *time.Time      `json:"published_at"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	User         UserPublic      `json:"user"`
	Qualities    []VideoQuality  `json:"qualities,omitempty"`
	Category     *Category       `json:"category,omitempty"`
}

type VideoQuality struct {
	ID          string `json:"id"`
	VideoID     string `json:"video_id"`
	Quality     string `json:"quality"`
	ManifestURL string `json:"manifest_url"`
	Bitrate     *int   `json:"bitrate"`
	FileSize    *int64 `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
}

var ErrVideoNotFound = errors.New("video not found")

type VideoListParams struct {
	Page       int
	Size       int
	CategoryID *int
	Sort       string // "latest", "trending"
	UserID     string // filter by uploader
	Status     string
}

func CreateVideo(ctx context.Context, db *sql.DB, userID string, title string, description string, categoryID *int, tags []string) (Video, error) {
	var v Video
	err := db.QueryRowContext(ctx, `
	INSERT INTO videos(user_id, title, description, category_id, tags)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, user_id, title, description, cover_url, duration, status, category_id, tags,
	          view_count, like_count, comment_count, share_count, published_at, created_at, updated_at
	`, userID, title, description, categoryID, pq.Array(tags)).Scan(
		&v.ID, &v.UserID, &v.Title, &v.Description, &v.CoverURL, &v.Duration, &v.Status,
		&v.CategoryID, pq.Array(&v.Tags),
		&v.ViewCount, &v.LikeCount, &v.CommentCount, &v.ShareCount,
		&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
	)
	return v, err
}

func GetVideoByID(ctx context.Context, db *sql.DB, id string) (Video, error) {
	var v Video
	err := db.QueryRowContext(ctx, `
	SELECT v.id, v.user_id, v.title, v.description, v.cover_url, v.duration, v.status,
	       v.category_id, v.tags, v.view_count, v.like_count, v.comment_count, v.share_count,
	       v.published_at, v.created_at, v.updated_at,
	       u.username, u.avatar_url
	FROM videos v JOIN users u ON v.user_id = u.id
	WHERE v.id = $1 AND v.status != 'deleted'
	`, id).Scan(
		&v.ID, &v.UserID, &v.Title, &v.Description, &v.CoverURL, &v.Duration, &v.Status,
		&v.CategoryID, pq.Array(&v.Tags),
		&v.ViewCount, &v.LikeCount, &v.CommentCount, &v.ShareCount,
		&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
		&v.Username, &v.AvatarURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Video{}, ErrVideoNotFound
		}
		return Video{}, err
	}
	return v, nil
}

func (v Video) Detail(qualities []VideoQuality, cat *Category) VideoDetail {
	d := VideoDetail{
		ID:           v.ID,
		UserID:       v.UserID,
		Title:        v.Title,
		Description:  v.Description,
		CoverURL:     v.CoverURL,
		Duration:     v.Duration,
		Status:       v.Status,
		CategoryID:   v.CategoryID,
		Tags:         v.Tags,
		ViewCount:    v.ViewCount,
		LikeCount:    v.LikeCount,
		CommentCount: v.CommentCount,
		ShareCount:   v.ShareCount,
		PublishedAt:  v.PublishedAt,
		CreatedAt:    v.CreatedAt,
		UpdatedAt:    v.UpdatedAt,
		User: UserPublic{
			ID:        v.UserID,
			Username:  v.Username,
			AvatarURL: v.AvatarURL,
		},
		Qualities: qualities,
		Category:  cat,
	}
	if d.Tags == nil {
		d.Tags = []string{}
	}
	if d.Qualities == nil {
		d.Qualities = []VideoQuality{}
	}
	return d
}

func ListVideos(ctx context.Context, db *sql.DB, params VideoListParams) ([]Video, int, error) {
	where := "WHERE v.status = 'published'"
	args := []any{}
	argIdx := 1

	if params.UserID != "" {
		where += " AND v.user_id = $" + itoa(argIdx)
		args = append(args, params.UserID)
		argIdx++
	}
	if params.CategoryID != nil {
		where += " AND v.category_id = $" + itoa(argIdx)
		args = append(args, *params.CategoryID)
		argIdx++
	}

	var total int
	countQuery := "SELECT count(*) FROM videos v " + where
	if err := db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	order := "v.published_at DESC"
	if params.Sort == "trending" {
		order = "v.view_count DESC"
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 || params.Size > 50 {
		params.Size = 20
	}
	offset := (params.Page - 1) * params.Size

	query := `
	SELECT v.id, v.user_id, v.title, v.description, v.cover_url, v.duration, v.status,
	       v.category_id, v.tags, v.view_count, v.like_count, v.comment_count, v.share_count,
	       v.published_at, v.created_at, v.updated_at,
	       u.username, u.avatar_url
	FROM videos v JOIN users u ON v.user_id = u.id
	` + where + ` ORDER BY ` + order + ` LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)
	args = append(args, params.Size, offset)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.CoverURL, &v.Duration, &v.Status,
			&v.CategoryID, pq.Array(&v.Tags),
			&v.ViewCount, &v.LikeCount, &v.CommentCount, &v.ShareCount,
			&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
			&v.Username, &v.AvatarURL,
		); err != nil {
			return nil, 0, err
		}
		videos = append(videos, v)
	}
	return videos, total, rows.Err()
}

func UpdateVideoMeta(ctx context.Context, db *sql.DB, id string, title string, description string, tags []string) (Video, error) {
	var v Video
	err := db.QueryRowContext(ctx, `
	UPDATE videos SET title = $2, description = $3, tags = $4, updated_at = now()
	WHERE id = $1 AND status != 'deleted'
	RETURNING id, user_id, title, description, cover_url, duration, status, category_id, tags,
	          view_count, like_count, comment_count, share_count, published_at, created_at, updated_at
	`, id, title, description, pq.Array(tags)).Scan(
		&v.ID, &v.UserID, &v.Title, &v.Description, &v.CoverURL, &v.Duration, &v.Status,
		&v.CategoryID, pq.Array(&v.Tags),
		&v.ViewCount, &v.LikeCount, &v.CommentCount, &v.ShareCount,
		&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Video{}, ErrVideoNotFound
		}
		return Video{}, err
	}
	return v, nil
}

func SoftDeleteVideo(ctx context.Context, db *sql.DB, id string) error {
	res, err := db.ExecContext(ctx, `
	UPDATE videos SET status = 'deleted', updated_at = now()
	WHERE id = $1 AND status != 'deleted'
	`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrVideoNotFound
	}
	return nil
}

func UpdateVideoStatus(ctx context.Context, db *sql.DB, id string, status string, duration int, coverURL string) error {
	_, err := db.ExecContext(ctx, `
	UPDATE videos SET status = $2, duration = $3, cover_url = $4, updated_at = now()
	WHERE id = $1
	`, id, status, duration, coverURL)
	return err
}

func PublishVideo(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(ctx, `
	UPDATE videos SET status = 'published', published_at = now(), updated_at = now()
	WHERE id = $1 AND status = 'processing'
	`, id)
	return err
}

func IncrementViewCount(ctx context.Context, db *sql.DB, id string) error {
	_, err := db.ExecContext(ctx, `
	UPDATE videos SET view_count = view_count + 1 WHERE id = $1
	`, id)
	return err
}

func GetVideoQualities(ctx context.Context, db *sql.DB, videoID string) ([]VideoQuality, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT id, video_id, quality, manifest_url, bitrate, file_size, created_at
	FROM video_qualities WHERE video_id = $1 ORDER BY bitrate DESC
	`, videoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quals []VideoQuality
	for rows.Next() {
		var q VideoQuality
		if err := rows.Scan(&q.ID, &q.VideoID, &q.Quality, &q.ManifestURL, &q.Bitrate, &q.FileSize, &q.CreatedAt); err != nil {
			return nil, err
		}
		quals = append(quals, q)
	}
	return quals, rows.Err()
}

func CreateVideoQuality(ctx context.Context, db *sql.DB, videoID string, quality string, manifestURL string, bitrate int, fileSize int64) (VideoQuality, error) {
	var q VideoQuality
	err := db.QueryRowContext(ctx, `
	INSERT INTO video_qualities(video_id, quality, manifest_url, bitrate, file_size)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, video_id, quality, manifest_url, bitrate, file_size, created_at
	`, videoID, quality, manifestURL, bitrate, fileSize).Scan(
		&q.ID, &q.VideoID, &q.Quality, &q.ManifestURL, &q.Bitrate, &q.FileSize, &q.CreatedAt,
	)
	return q, err
}

func itoa(n int) string {
	return fmt.Sprintf("%d", n)
}
