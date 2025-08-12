package config

import "github.com/s21platform/jarvis-bot/internal/model"

// ProjectMapping содержит маппинг каналов Mattermost в проекты Jira
type ProjectMapping struct {
	// ChannelToProject маппит название канала в ключ проекта Jira
	ChannelToProject map[string]string
}

// NewProjectMapping создает новый маппинг проектов
func NewProjectMapping() *ProjectMapping {
	// TODO: В будущем можно загружать из файла конфигурации или БД
	return &ProjectMapping{
		ChannelToProject: map[string]string{
			// Пример маппинга:
			// "chat-team-backend": "BACK",
			// "chat-team-frontend": "FRONT",

			// Community
			model.COMMUNITY_SERVICE_CHANNEL:        model.COMMUNITY_PROJECT,
			model.COMMUNITY_SERVICE_NEWS_CHANNEL:   model.COMMUNITY_PROJECT,
			model.COMMUNITY_SERVICE_PUBLIC_CHANNEL: model.COMMUNITY_PROJECT,

			// Materials
			model.MATERIALS_SERVICE_CHANNEL:        model.MATERIALS_PROJECT,
			model.MATERIALS_SERVICE_NEWS_CHANNEL:   model.MATERIALS_PROJECT,
			model.MATERIALS_SERVICE_PUBLIC_CHANNEL: model.MATERIALS_PROJECT,

			// Jarvis
			model.JARVIS_BOT_CHANNEL:      model.JARVIS_PROJECT,
			model.JARVIS_BOT_NEWS_CHANNEL: model.JARVIS_PROJECT,
			model.JARVIS_PUBLIC_CHANNEL:   model.JARVIS_PROJECT,

			// Frontend
			model.FRONTEND_CHANNEL:       model.FRONTEND_PROJECT,
			model.FRONTEND_NEWS_CHANNEL:  model.FRONTEND_PROJECT,
			model.FRONTEND_ULTRA_CHANNEL: model.FRONTEND_PROJECT,

			// Tech Stream
			model.TECH_STREAM_CHANNEL: model.TECH_STREAM,

			// Evo
			model.EVO_CHANNEL: model.EVO_PROJECT,
		},
	}
}

// GetProjectKey возвращает ключ проекта Jira для канала
func (p *ProjectMapping) GetProjectKey(channelName string) string {
	if projectKey, ok := p.ChannelToProject[channelName]; ok {
		return projectKey
	}
	return "" // Возвращаем пустую строку, если маппинг не найден
}
