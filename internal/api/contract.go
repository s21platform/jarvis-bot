package api

import (
	"context"

	"github.com/s21platform/jarvis-bot/internal/model"
)

type DatabaseRepo interface {
	GetBirthday(ctx context.Context, userId string) (*model.Birthday, error)
	SetBirthday(ctx context.Context, birthday *model.Birthday, userId string, channelId string, nickname string) error
}
