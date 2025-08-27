package worker

import (
	"context"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	models "github.com/s21platform/jarvis-bot/internal/model"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type MattermostClient interface {
	GetUsers(page int, perPage int, etag string) ([]*model.User, *model.Response, error)
	GetMe(etag string) (*model.User, *model.Response, error)
	CreateDirectChannel(userId1 string, userId2 string) (*model.Channel, *model.Response, error)
	CreatePost(post *model.Post) (*model.Post, *model.Response, error)
}

type DbRepo interface {
	GetConfig(ctx context.Context, key string) (models.Config, error)
	GetBirthday(ctx context.Context, userID string) (*models.Birthday, error)
	GetAllBirthdays(ctx context.Context) (map[string]time.Time, error)
}

type Metrics interface {
	Increment(name string)
	Count(name string, value int64)
	Duration(startTime int64, name string)
}
