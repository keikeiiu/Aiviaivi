package model

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrAlreadyExists = errors.New("already exists")
var ErrNotExists = errors.New("not exists")

// --- Likes ---

func LikeVideo(ctx context.Context, db *sql.DB, userID string, videoID string) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO likes(user_id, video_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, userID, videoID)
	if err != nil {
		return err
	}
	_, _ = db.ExecContext(ctx, `UPDATE videos SET like_count = (SELECT count(*) FROM likes WHERE video_id = $1) WHERE id = $1`, videoID)
	return nil
}

func UnlikeVideo(ctx context.Context, db *sql.DB, userID string, videoID string) error {
	res, err := db.ExecContext(ctx, `DELETE FROM likes WHERE user_id = $1 AND video_id = $2`, userID, videoID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotExists
	}
	_, _ = db.ExecContext(ctx, `UPDATE videos SET like_count = (SELECT count(*) FROM likes WHERE video_id = $1) WHERE id = $1`, videoID)
	return nil
}

func HasUserLiked(ctx context.Context, db *sql.DB, userID string, videoID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND video_id = $2)`, userID, videoID,
	).Scan(&exists)
	return exists, err
}

// --- Favorites ---

func FavoriteVideo(ctx context.Context, db *sql.DB, userID string, videoID string) error {
	_, err := db.ExecContext(ctx, `
	INSERT INTO favorites(user_id, video_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, userID, videoID)
	return err
}

func UnfavoriteVideo(ctx context.Context, db *sql.DB, userID string, videoID string) error {
	res, err := db.ExecContext(ctx, `DELETE FROM favorites WHERE user_id = $1 AND video_id = $2`, userID, videoID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotExists
	}
	return nil
}

func HasUserFavorited(ctx context.Context, db *sql.DB, userID string, videoID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM favorites WHERE user_id = $1 AND video_id = $2)`, userID, videoID,
	).Scan(&exists)
	return exists, err
}

func ListFavorites(ctx context.Context, db *sql.DB, userID string, page int, size int) ([]Video, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	offset := (page - 1) * size

	var total int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM favorites WHERE user_id = $1`, userID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := db.QueryContext(ctx, `
	SELECT v.id, v.user_id, v.title, v.description, v.cover_url, v.duration, v.status,
	       v.category_id, v.tags, v.view_count, v.like_count, v.comment_count, v.share_count,
	       v.published_at, v.created_at, v.updated_at,
	       u.username, u.avatar_url
	FROM favorites f
	JOIN videos v ON f.video_id = v.id
	JOIN users u ON v.user_id = u.id
	WHERE f.user_id = $1 AND v.status = 'published'
	ORDER BY f.created_at DESC
	LIMIT $2 OFFSET $3
	`, userID, size, offset)
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

// --- Follows ---

func FollowUser(ctx context.Context, db *sql.DB, followerID string, creatorID string) error {
	if followerID == creatorID {
		return errors.New("cannot follow yourself")
	}
	_, err := db.ExecContext(ctx, `
	INSERT INTO follows(follower_id, creator_id) VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, followerID, creatorID)
	return err
}

func UnfollowUser(ctx context.Context, db *sql.DB, followerID string, creatorID string) error {
	res, err := db.ExecContext(ctx,
		`DELETE FROM follows WHERE follower_id = $1 AND creator_id = $2`, followerID, creatorID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotExists
	}
	return nil
}

func IsFollowing(ctx context.Context, db *sql.DB, followerID string, creatorID string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx,
		`SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND creator_id = $2)`, followerID, creatorID,
	).Scan(&exists)
	return exists, err
}

func GetFollowerCount(ctx context.Context, db *sql.DB, userID string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT count(*) FROM follows WHERE creator_id = $1`, userID).Scan(&count)
	return count, err
}

func GetFollowingCount(ctx context.Context, db *sql.DB, userID string) (int, error) {
	var count int
	err := db.QueryRowContext(ctx, `SELECT count(*) FROM follows WHERE follower_id = $1`, userID).Scan(&count)
	return count, err
}

// --- Search ---

func SearchVideos(ctx context.Context, db *sql.DB, query string, page int, size int) ([]Video, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	offset := (page - 1) * size

	var total int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM videos WHERE status = 'published' AND search_vector @@ plainto_tsquery('english', $1)`,
		query,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := db.QueryContext(ctx, `
	SELECT v.id, v.user_id, v.title, v.description, v.cover_url, v.duration, v.status,
	       v.category_id, v.tags, v.view_count, v.like_count, v.comment_count, v.share_count,
	       v.published_at, v.created_at, v.updated_at,
	       u.username, u.avatar_url,
	       ts_rank(v.search_vector, plainto_tsquery('english', $1)) AS rank
	FROM videos v JOIN users u ON v.user_id = u.id
	WHERE v.status = 'published' AND v.search_vector @@ plainto_tsquery('english', $1)
	ORDER BY rank DESC, v.view_count DESC
	LIMIT $2 OFFSET $3
	`, query, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var videos []Video
	for rows.Next() {
		var v Video
		var rank float64
		if err := rows.Scan(
			&v.ID, &v.UserID, &v.Title, &v.Description, &v.CoverURL, &v.Duration, &v.Status,
			&v.CategoryID, pq.Array(&v.Tags),
			&v.ViewCount, &v.LikeCount, &v.CommentCount, &v.ShareCount,
			&v.PublishedAt, &v.CreatedAt, &v.UpdatedAt,
			&v.Username, &v.AvatarURL,
			&rank,
		); err != nil {
			return nil, 0, err
		}
		videos = append(videos, v)
	}
	return videos, total, rows.Err()
}

func SearchUsers(ctx context.Context, db *sql.DB, query string, page int, size int) ([]User, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	offset := (page - 1) * size

	searchTerm := "%" + query + "%"

	var total int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM users WHERE username ILIKE $1 OR bio ILIKE $1`, searchTerm,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := db.QueryContext(ctx, `
	SELECT id, username, email, password_hash, avatar_url, bio, role, created_at, updated_at
	FROM users WHERE username ILIKE $1 OR bio ILIKE $1
	ORDER BY username ASC
	LIMIT $2 OFFSET $3
	`, searchTerm, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.PasswordHash,
			&u.AvatarURL, &u.Bio, &u.Role, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

