package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/model"
)

type Postgres struct {
	conn *sqlx.DB
}

func New(cfg *config.Config) *Postgres {
	conStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Database)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	return &Postgres{conn: conn}
}

func (p *Postgres) Close() {
	_ = p.conn.Close()
}

func (p *Postgres) GetCred(ctx context.Context, name string) (model.Data, error) {
	query, args, err := sq.Select(`data`).
		From(`credentials`).
		Where(sq.Eq{"service": name}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return model.Data{}, fmt.Errorf("failed to build query: %w", err)
	}

	var credData model.CredData
	err = p.conn.GetContext(ctx, &credData, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Data{}, sql.ErrNoRows
		}
		return model.Data{}, fmt.Errorf("failed to get credentials from db: %w", err)
	}

	var data model.Data
	err = json.Unmarshal(credData.Data, &data)
	if err != nil {
		return model.Data{}, fmt.Errorf("failed to unmarshal cred: %w", err)
	}

	return data, nil
}

// AddNews добавляет новость в базу данных
func (p *Postgres) AddNews(ctx context.Context, threadID, channel, messageID string) error {
	query, args, err := sq.Insert("news").
		Columns("thread_id", "channel", "message_id").
		Values(threadID, channel, messageID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = p.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert news: %w", err)
	}

	return nil
}

// DeleteNews удаляет новость из базы данных
func (p *Postgres) DeleteNews(ctx context.Context, threadID string) error {
	query, args, err := sq.Delete("news").
		Where(sq.Eq{"thread_id": threadID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	result, err := p.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete news: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetNews возвращает список всех новостей за последние 3 недели
func (p *Postgres) GetNews(ctx context.Context) ([]model.News, error) {
	query, args, err := sq.Select("id", "thread_id", "channel", "message_id", "created_at").
		From("news").
		Where("created_at > NOW() - INTERVAL '3 weeks'").
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var news []model.News
	err = p.conn.SelectContext(ctx, &news, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get news: %w", err)
	}

	return news, nil
}

// IsNewsExists проверяет существование новости в базе данных
func (p *Postgres) IsNewsExists(ctx context.Context, threadID string) (bool, error) {
	query, args, err := sq.Select("1").
		From("news").
		Where(sq.Eq{"thread_id": threadID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = p.conn.GetContext(ctx, &exists, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check news existence: %w", err)
	}

	return exists, nil
}
