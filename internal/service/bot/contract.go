package bot

import (
	"context"
	"github.com/s21platform/jarvis-bot/internal/model"
)

type DbRepo interface {
	GetCred(ctx context.Context, name string) (model.Data, error)
}
