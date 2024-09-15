package postgres

import (
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

func (p *Postgres) CreateTask(channelName, taskType, taskTitle, assignee string) (int64, error) {
	query, args, err := sq.Insert("tasks").
		Columns("service", "task_type", "task_title", "assignee").
		Values(channelName, taskType, taskTitle, assignee).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, err
	}

	var insertedID int64
	err = p.conn.QueryRow(query, args...).Scan(&insertedID)
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}

func (p *Postgres) GetTasksByUUID(assignee, service string) ([]model.TasksByUUID, error) {
	query, args, err := sq.Select(
		"id",
		`task_type`,
		`task_title`,
		`task_description`,
	).
		From("tasks").
		Where(sq.And{
			sq.Eq{"assignee": assignee},
			sq.Eq{"service": service},
		}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	var tasksByUUID []model.TasksByUUID
	err = p.conn.Select(&tasksByUUID, query, args...)
	if err != nil {
		return nil, err
	}
	return tasksByUUID, nil
}
