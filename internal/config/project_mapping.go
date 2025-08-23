package config

import "github.com/s21platform/jarvis-bot/internal/model"

// ProjectConfig содержит конфигурацию проекта Jira для канала
type ProjectConfig struct {
	// ProjectKey ключ проекта в Jira
	ProjectKey string `json:"project_key"`
	// Label опциональная метка для задач из этого канала
	Label *string `json:"label,omitempty"`
}

// ProjectMapping содержит маппинг каналов Mattermost в проекты Jira
type ProjectMapping struct {
	// ChannelToProject маппит название канала в конфигурацию проекта Jira
	ChannelToProject map[string]ProjectConfig
}

// NewProjectMapping создает новый маппинг проектов
func NewProjectMapping() *ProjectMapping {
	// TODO: В будущем можно загружать из файла конфигурации или БД
	return &ProjectMapping{
		ChannelToProject: map[string]ProjectConfig{

			// Community
			model.COMMUNITY_SERVICE_CHANNEL: {
				ProjectKey: model.COMMUNITY_PROJECT,
			},
			model.COMMUNITY_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.COMMUNITY_PROJECT,
			},
			model.COMMUNITY_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.COMMUNITY_PROJECT,
			},

			// Materials
			model.MATERIALS_SERVICE_CHANNEL: {
				ProjectKey: model.MATERIALS_PROJECT,
			},
			model.MATERIALS_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.MATERIALS_PROJECT,
			},
			model.MATERIALS_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.MATERIALS_PROJECT,
			},

			// Jarvis
			model.JARVIS_BOT_CHANNEL: {
				ProjectKey: model.JARVIS_PROJECT,
			},
			model.JARVIS_BOT_NEWS_CHANNEL: {
				ProjectKey: model.JARVIS_PROJECT,
			},
			model.JARVIS_PUBLIC_CHANNEL: {
				ProjectKey: model.JARVIS_PROJECT,
			},

			// Frontend
			model.FRONTEND_CHANNEL: {
				ProjectKey: model.FRONTEND_PROJECT,
			},
			model.FRONTEND_NEWS_CHANNEL: {
				ProjectKey: model.FRONTEND_PROJECT,
			},
			model.FRONTEND_ULTRA_CHANNEL: {
				ProjectKey: model.FRONTEND_PROJECT,
			},

			// Tech Stream
			model.TECH_STREAM_CHANNEL: {
				ProjectKey: model.TECH_STREAM,
			},

			// Evo
			model.EVO_CHANNEL: {
				ProjectKey: model.EVO_PROJECT,
			},

			// User
			model.USER_SERVICE_CHANNEL: {
				ProjectKey: model.USER_PROJECT,
			},
			model.USER_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.USER_PROJECT,
			},
			model.USER_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.USER_PROJECT,
			},

			// Storage
			model.STORAGE_SERVICE_CHANNEL: {
				ProjectKey: model.STORAGE_PROJECT,
			},
			model.STORAGE_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.STORAGE_PROJECT,
			},
			model.STORAGE_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.STORAGE_PROJECT,
			},

			// Advert
			model.ADVERT_SERVICE_CHANNEL: {
				ProjectKey: model.ADVERT_PROJECT,
			},
			model.ADVERT_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.ADVERT_PROJECT,
			},
			model.ADVERT_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.ADVERT_PROJECT,
			},

			// Auth
			model.AUTH_SERVICE_CHANNEL: {
				ProjectKey: model.AUTH_PROJECT,
			},
			model.AUTH_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.AUTH_PROJECT,
			},
			model.AUTH_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.AUTH_PROJECT,
			},

			// Feed
			model.FEED_SERVICE_CHANNEL: {
				ProjectKey: model.FEED_PROJECT,
			},
			model.FEED_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.FEED_PROJECT,
			},
			model.FEED_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.FEED_PROJECT,
			},

			// Search
			model.SEARCH_SERVICE_CHANNEL: {
				ProjectKey: model.SEARCH_PROJECT,
			},
			model.SEARCH_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.SEARCH_PROJECT,
			},
			model.SEARCH_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.SEARCH_PROJECT,
			},

			// Society
			model.SOCIETY_SERVICE_CHANNEL: {
				ProjectKey: model.SOCIETY_PROJECT,
			},
			model.SOCIETY_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.SOCIETY_PROJECT,
			},
			model.SOCIETY_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.SOCIETY_PROJECT,
			},

			// Optionhub
			model.OPTIONHUB_SERVICE_CHANNEL: {
				ProjectKey: model.OPTIONHUB_PROJECT,
			},
			model.OPTIONHUB_SERVICE_NEWS_CHANNEL: {
				ProjectKey: model.OPTIONHUB_PROJECT,
			},
			model.OPTIONHUB_SERVICE_PUBLIC_CHANNEL: {
				ProjectKey: model.OPTIONHUB_PROJECT,
			},

			// Metrics
			model.LOGGER_LIB_CHANNEL: {
				ProjectKey: model.TECH_STREAM,
				Label:      model.LOGGER_LABEL,
			},
		},
	}
}

// GetProjectKey возвращает ключ проекта Jira для канала
func (p *ProjectMapping) GetProjectKey(channelName string) string {
	if config, ok := p.ChannelToProject[channelName]; ok {
		return config.ProjectKey
	}
	return "" // Возвращаем пустую строку, если маппинг не найден
}

// GetLabel возвращает метку для задач из канала, если она установлена
func (p *ProjectMapping) GetLabel(channelName string) *string {
	if config, ok := p.ChannelToProject[channelName]; ok {
		return config.Label
	}
	return nil
}
