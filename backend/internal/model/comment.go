package model

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

type Comment struct {
	ID        string    `json:"id"`
	VideoID   string    `json:"video_id"`
	UserID    string    `json:"user_id"`
	ParentID  *string   `json:"parent_id"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Joined
	Username  string     `json:"username"`
	AvatarURL string     `json:"avatar_url"`
	Replies   []Comment  `json:"replies,omitempty"`
}

var ErrCommentNotFound = errors.New("comment not found")

type CommentListParams struct {
	Page   int
	Size   int
	Sort   string // "latest", "hot"
}

func CreateComment(ctx context.Context, db *sql.DB, videoID string, userID string, content string, parentID *string) (Comment, error) {
	var c Comment
	err := db.QueryRowContext(ctx, `
	INSERT INTO comments(video_id, user_id, content, parent_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, video_id, user_id, parent_id, content, like_count, created_at, updated_at
	`, videoID, userID, content, parentID).Scan(
		&c.ID, &c.VideoID, &c.UserID, &c.ParentID, &c.Content,
		&c.LikeCount, &c.CreatedAt, &c.UpdatedAt,
	)
	return c, err
}

func ListComments(ctx context.Context, db *sql.DB, videoID string, params CommentListParams) ([]Comment, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Size < 1 || params.Size > 50 {
		params.Size = 20
	}
	offset := (params.Page - 1) * params.Size

	order := "c.created_at DESC"
	if params.Sort == "hot" {
		order = "c.like_count DESC, c.created_at DESC"
	}

	// Count top-level comments only
	var total int
	if err := db.QueryRowContext(ctx,
		`SELECT count(*) FROM comments WHERE video_id = $1 AND parent_id IS NULL`, videoID,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Get top-level comments
	rows, err := db.QueryContext(ctx, `
	SELECT c.id, c.video_id, c.user_id, c.parent_id, c.content, c.like_count, c.created_at, c.updated_at,
	       u.username, u.avatar_url
	FROM comments c JOIN users u ON c.user_id = u.id
	WHERE c.video_id = $1 AND c.parent_id IS NULL
	ORDER BY `+order+` LIMIT $2 OFFSET $3
	`, videoID, params.Size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var comments []Comment
	var commentIDs []string
	for rows.Next() {
		var c Comment
		if err := rows.Scan(
			&c.ID, &c.VideoID, &c.UserID, &c.ParentID, &c.Content,
			&c.LikeCount, &c.CreatedAt, &c.UpdatedAt,
			&c.Username, &c.AvatarURL,
		); err != nil {
			return nil, 0, err
		}
		c.Replies = []Comment{}
		comments = append(comments, c)
		commentIDs = append(commentIDs, c.ID)
	}

	// Fetch replies for top-level comments
	if len(commentIDs) > 0 {
		replies, err := getReplies(ctx, db, commentIDs)
		if err != nil {
			return nil, 0, err
		}
		replyMap := map[string][]Comment{}
		for _, r := range replies {
			if r.ParentID != nil {
				replyMap[*r.ParentID] = append(replyMap[*r.ParentID], r)
			}
		}
		for i := range comments {
			if reps, ok := replyMap[comments[i].ID]; ok {
				comments[i].Replies = reps
			}
		}
	}

	return comments, total, nil
}

func getReplies(ctx context.Context, db *sql.DB, parentIDs []string) ([]Comment, error) {
	query := `
	SELECT c.id, c.video_id, c.user_id, c.parent_id, c.content, c.like_count, c.created_at, c.updated_at,
	       u.username, u.avatar_url
	FROM comments c JOIN users u ON c.user_id = u.id
	WHERE c.parent_id = ANY($1::uuid[])
	ORDER BY c.created_at ASC
	`
	rows, err := db.QueryContext(ctx, query, pq.Array(parentIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(
			&c.ID, &c.VideoID, &c.UserID, &c.ParentID, &c.Content,
			&c.LikeCount, &c.CreatedAt, &c.UpdatedAt,
			&c.Username, &c.AvatarURL,
		); err != nil {
			return nil, err
		}
		c.Replies = []Comment{}
		replies = append(replies, c)
	}
	return replies, rows.Err()
}

func DeleteComment(ctx context.Context, db *sql.DB, id string, userID string) error {
	res, err := db.ExecContext(ctx,
		`DELETE FROM comments WHERE id = $1 AND user_id = $2`, id, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrCommentNotFound
	}
	return nil
}
