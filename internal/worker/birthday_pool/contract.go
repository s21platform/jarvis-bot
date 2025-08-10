package worker

import (
	"context"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	modelInternal "github.com/s21platform/jarvis-bot/internal/model"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type MattermostClient interface {
	GetMe(etag string) (*model.User, *model.Response, error)
	CreateDirectChannel(userId1 string, userId2 string) (*model.Channel, *model.Response, error)
	CreatePost(post *model.Post) (*model.Post, *model.Response, error)
	GetUsers(page int, perPage int, etag string) ([]*model.User, *model.Response, error)
}

type DbRepo interface {
	GetBirthday(ctx context.Context, userId string) (*modelInternal.Birthday, error)
	SetBirthday(ctx context.Context, birthday *modelInternal.Birthday, userId string, channelId string, nickname string) error
	GetConfig(ctx context.Context, key string) (modelInternal.Config, error)
}

type Metrics interface {
	Count(name string, value int64)
	Disconnect()
	Duration(timestamp int64, name string)
	Gauge(name string, value float64)
	Increment(name string)
}
