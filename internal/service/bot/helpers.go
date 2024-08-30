package bot

import (
	"encoding/json"
	"errors"
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

func parseCommand(message string) (*Command, error) {
	parts := strings.Fields(strings.Replace(message, "@jarvis", "", 1))
	if len(parts) < 2 {
		return nil, errors.New("invalid command")
	}
	return &Command{
		Name: parts[0],
		Cmd:  strings.Join(parts[1:], " "),
	}, nil
}
