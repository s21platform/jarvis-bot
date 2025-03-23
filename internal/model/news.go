package model

import "time"

// News представляет собой модель новости
type News struct {
	ID        int64     `db:"id"`
	ThreadID  string    `db:"thread_id"`
	Channel   string    `db:"channel"`
	MessageID string    `db:"message_id"`
	CreatedAt time.Time `db:"created_at"`
}
