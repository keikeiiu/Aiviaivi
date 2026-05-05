package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Playlist struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	VideoCount  int       `json:"video_count"`
	// Joined
	Username string `json:"username,omitempty"`
}

type PlaylistVideo struct {
	PlaylistID string    `json:"playlist_id"`
	VideoID    string    `json:"video_id"`
	Position   int       `json:"position"`
	AddedAt    time.Time `json:"added_at"`
	// Joined video fields
	Title    string `json:"title"`
	CoverURL string `json:"cover_url"`
	Duration int    `json:"duration"`
}

var ErrPlaylistNotFound = errors.New("playlist not found")
var ErrPlaylistVideoNotFound = errors.New("video not in playlist")

func CreatePlaylist(ctx context.Context, db *sql.DB, userID string, name string, description string, isPublic bool) (Playlist, error) {
	var p Playlist
	err := db.QueryRowContext(ctx, `
	INSERT INTO playlists(user_id, name, description, is_public)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, name, description, is_public, created_at, updated_at
	`, userID, name, description, isPublic).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Description, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt,
	)
	return p, err
}

func GetPlaylistByID(ctx context.Context, db *sql.DB, id string) (Playlist, error) {
	var p Playlist
	err := db.QueryRowContext(ctx, `
	SELECT p.id, p.user_id, p.name, p.description, p.is_public, p.created_at, p.updated_at,
	       COALESCE((SELECT count(*) FROM playlist_videos WHERE playlist_id = p.id), 0),
	       u.username
	FROM playlists p JOIN users u ON p.user_id = u.id
	WHERE p.id = $1
	`, id).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Description, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt,
		&p.VideoCount, &p.Username,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Playlist{}, ErrPlaylistNotFound
		}
		return Playlist{}, err
	}
	return p, nil
}

func ListUserPlaylists(ctx context.Context, db *sql.DB, userID string) ([]Playlist, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT p.id, p.user_id, p.name, p.description, p.is_public, p.created_at, p.updated_at,
	       COALESCE((SELECT count(*) FROM playlist_videos WHERE playlist_id = p.id), 0)
	FROM playlists p
	WHERE p.user_id = $1
	ORDER BY p.updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var p Playlist
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Description, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt, &p.VideoCount); err != nil {
			return nil, err
		}
		playlists = append(playlists, p)
	}
	return playlists, rows.Err()
}

func UpdatePlaylist(ctx context.Context, db *sql.DB, id string, name string, description string, isPublic bool) (Playlist, error) {
	var p Playlist
	err := db.QueryRowContext(ctx, `
	UPDATE playlists SET name = $2, description = $3, is_public = $4, updated_at = now()
	WHERE id = $1
	RETURNING id, user_id, name, description, is_public, created_at, updated_at
	`, id, name, description, isPublic).Scan(
		&p.ID, &p.UserID, &p.Name, &p.Description, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Playlist{}, ErrPlaylistNotFound
		}
		return Playlist{}, err
	}
	return p, nil
}

func DeletePlaylist(ctx context.Context, db *sql.DB, id string) error {
	res, err := db.ExecContext(ctx, `DELETE FROM playlists WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPlaylistNotFound
	}
	return nil
}

func AddVideoToPlaylist(ctx context.Context, db *sql.DB, playlistID string, videoID string) error {
	// Get next position
	var maxPos int
	db.QueryRowContext(ctx,
		`SELECT COALESCE(MAX(position), -1) FROM playlist_videos WHERE playlist_id = $1`, playlistID,
	).Scan(&maxPos)

	_, err := db.ExecContext(ctx, `
	INSERT INTO playlist_videos(playlist_id, video_id, position)
	VALUES ($1, $2, $3)
	ON CONFLICT (playlist_id, video_id) DO NOTHING
	`, playlistID, videoID, maxPos+1)
	return err
}

func RemoveVideoFromPlaylist(ctx context.Context, db *sql.DB, playlistID string, videoID string) error {
	res, err := db.ExecContext(ctx,
		`DELETE FROM playlist_videos WHERE playlist_id = $1 AND video_id = $2`, playlistID, videoID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrPlaylistVideoNotFound
	}
	return nil
}

func GetPlaylistVideos(ctx context.Context, db *sql.DB, playlistID string) ([]PlaylistVideo, error) {
	rows, err := db.QueryContext(ctx, `
	SELECT pv.playlist_id, pv.video_id, pv.position, pv.added_at,
	       v.title, v.cover_url, v.duration
	FROM playlist_videos pv
	JOIN videos v ON pv.video_id = v.id
	WHERE pv.playlist_id = $1
	ORDER BY pv.position ASC
	`, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []PlaylistVideo
	for rows.Next() {
		var pv PlaylistVideo
		if err := rows.Scan(&pv.PlaylistID, &pv.VideoID, &pv.Position, &pv.AddedAt, &pv.Title, &pv.CoverURL, &pv.Duration); err != nil {
			return nil, err
		}
		videos = append(videos, pv)
	}
	return videos, rows.Err()
}
