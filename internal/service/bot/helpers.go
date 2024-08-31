package bot

import (
	"encoding/json"
	"github.com/mattermost/mattermost-server/v6/model"
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
