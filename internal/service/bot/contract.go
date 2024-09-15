package bot

import "github.com/s21platform/jarvis-bot/internal/model"

type DbRepo interface {
	CreateTask(channelName, taskType, taskTitle, assignee string) (int64, error)
	GetTasksByUUID(assignee, service string) ([]model.TasksByUUID, error)
}
