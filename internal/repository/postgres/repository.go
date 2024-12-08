package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/model"
	"log"
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
		return model.Data{}, fmt.Errorf("failed to get credentials from db: %w", err)
	}

	var data model.Data
	err = json.Unmarshal(credData.Data, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}
