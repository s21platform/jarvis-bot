package bot

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v6/model"
	servicemodel "github.com/s21platform/jarvis-bot/internal/model"
	"strconv"
	"strings"
)

type Command struct {
	Name string
	Cmd  string
}

func getPost(event *model.WebSocketEvent) (*model.Post, error) {
	jsonPost := event.GetData()["post"].(string)
	post := &model.Post{}
	err := json.Unmarshal([]byte(jsonPost), post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func parseCommand(message string) *Command {
	parts := strings.Fields(strings.Replace(message, "@jarvis", "", 1))
	var cmd string
	var name string
	switch len(parts) {
	case 0:
		name = "help"
		cmd = ""
	case 1:
		name = parts[0]
		cmd = ""
	default:
		name = parts[0]
		cmd = strings.Join(parts[1:], " ")
	}
	return &Command{
		Name: name,
		Cmd:  cmd,
	}
}

func CreateTable(headers []string, rows [][]string) string {
	var builder strings.Builder

	// Создаем строку заголовков
	builder.WriteString("|")
	for _, header := range headers {
		builder.WriteString(" " + header + " |")
	}
	builder.WriteString("\n")

	// Создаем разделитель
	builder.WriteString("|")
	for range headers {
		builder.WriteString("------------|")
	}
	builder.WriteString("\n")

	// Добавляем строки данных
	for _, row := range rows {
		builder.WriteString("|")
		for _, cell := range row {
			builder.WriteString(" " + cell + " |")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func convertModelToString(t []servicemodel.TasksByUUID) [][]string {
	result := make([][]string, len(t))
	for i, val := range t {
		description := ""
		if val.TaskDescription != nil {
			description = strings.Replace(*val.TaskDescription, "\\n", " ", -1)
		}
		result[i] = []string{strconv.FormatInt(val.ID, 10), val.TaskTitle, description, val.TaskType}
	}
	return result
}

func convertModelAllTasksToString(t []servicemodel.TasksByChannel) [][]string {
	result := make([][]string, len(t))
	for i, val := range t {
		description := ""
		if val.TaskDescription != nil {
			description = strings.Replace(*val.TaskDescription, "\\n", " ", -1)
		}
		result[i] = []string{strconv.FormatInt(val.ID, 10), val.Assignee, val.TaskTitle, description, val.TaskType}
	}
	return result
}
