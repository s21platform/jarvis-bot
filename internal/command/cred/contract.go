package cred

import (
	"context"

	"github.com/s21platform/jarvis-bot/internal/model"
)

type DbRepo interface {
	GetCred(ctx context.Context, key string) (model.Data, error)
}
