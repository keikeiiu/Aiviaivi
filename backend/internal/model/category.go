package model

import (
	"context"
	"database/sql"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ListCategories(ctx context.Context, db *sql.DB) ([]Category, error) {
	rows, err := db.QueryContext(ctx, `SELECT id, name, slug FROM categories ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

func GetCategoryByID(ctx context.Context, db *sql.DB, id int) (*Category, error) {
	var c Category
	err := db.QueryRowContext(ctx,
		`SELECT id, name, slug FROM categories WHERE id = $1`, id,
	).Scan(&c.ID, &c.Name, &c.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

var ErrCategoryNotFound = sql.ErrNoRows
