package create

import (
	"context"

	"github.com/s21platform/jarvis-bot/internal/model"
)

type DbRepo interface {
	GetConfig(ctx context.Context, key string) (model.Config, error)
}

type JiraClient interface {
	CreateIssue(title string, issueType string, projectKey string) (string, error)
	GetBaseURL() string
}
