package accents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Repo struct {
	database *sql.DB
}

func NewRepo(database *sql.DB) *Repo {
	return &Repo{database: database}
}

func (r *Repo) List(ctx context.Context) ([]Preset, error) {
	rows, err := r.database.QueryContext(ctx, `
		SELECT id, name, accent, accent_strong, accent_ink
		FROM accent_presets
		ORDER BY name COLLATE NOCASE, id
	`)
	if err != nil {
		return nil, fmt.Errorf("list accent presets: %w", err)
	}
	defer rows.Close()

	presets := []Preset{}
	for rows.Next() {
		var preset Preset
		if err := rows.Scan(&preset.ID, &preset.Name, &preset.Accent, &preset.Strong, &preset.Ink); err != nil {
			return nil, fmt.Errorf("scan accent preset: %w", err)
		}
		presets = append(presets, preset)
	}
	return presets, rows.Err()
}

func (r *Repo) Get(ctx context.Context, id int64) (Preset, error) {
	var preset Preset
	err := r.database.QueryRowContext(ctx, `
		SELECT id, name, accent, accent_strong, accent_ink
		FROM accent_presets
		WHERE id = ?
	`, id).Scan(&preset.ID, &preset.Name, &preset.Accent, &preset.Strong, &preset.Ink)
	if errors.Is(err, sql.ErrNoRows) {
		return Preset{}, ErrNotFound
	}
	if err != nil {
		return Preset{}, fmt.Errorf("get accent preset: %w", err)
	}
	return preset, nil
}

func (r *Repo) Create(ctx context.Context, preset Preset) (Preset, error) {
	result, err := r.database.ExecContext(ctx, `
		INSERT INTO accent_presets (name, accent, accent_strong, accent_ink)
		VALUES (?, ?, ?, ?)
	`, preset.Name, preset.Accent, preset.Strong, preset.Ink)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return Preset{}, ErrNameTaken
		}
		return Preset{}, fmt.Errorf("create accent preset: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Preset{}, fmt.Errorf("created accent preset id: %w", err)
	}
	return r.Get(ctx, id)
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	result, err := r.database.ExecContext(ctx, `DELETE FROM accent_presets WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete accent preset: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete accent preset rows: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
