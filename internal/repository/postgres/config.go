package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/s21platform/jarvis-bot/internal/model"
)

func (p *Postgres) GetConfig(ctx context.Context, key string) (model.Config, error) {
	query, args, err := sq.Select("value").
		From("config").
		Where(sq.Eq{"key": key}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return model.Config{}, fmt.Errorf("failed to build query: %w", err)
	}

	var config model.Config
	err = p.conn.GetContext(ctx, &config, query, args...)
	if err != nil {
		return model.Config{}, fmt.Errorf("failed to get config: %w", err)
	}

	return config, nil
}
