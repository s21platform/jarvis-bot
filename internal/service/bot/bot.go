package bot

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"log"
	"strings"
)

type Bot struct {
	websocket *model.WebSocketClient
	client    *model.Client4
	user      *model.User
}

func New(cfg *config.Config) *Bot {
	client := model.NewAPIv4Client(cfg.Url)
	client.SetOAuthToken(cfg.Token)

	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Не удалось получить информацию о пользователе: %v", err)
	}
	//log.Printf("Успешно авторизован как %s", user.Username)

	websocketClient, err := model.NewWebSocketClient4("wss://"+cfg.Url[len("https://"):], client.AuthToken)
	if err != nil {
		log.Fatalf("Ошибка подключения к WebSocket: %v", err)
	}

	return &Bot{
		websocket: websocketClient,
		client:    client,
		user:      user,
	}
}

func (b *Bot) Listen() {
	go func() {
		//time.Sleep(2 * time.Second)
		for event := range b.websocket.EventChannel {
			if event.EventType() == model.WebsocketEventPosted {
				post, err := getPost(event)
				if err != nil {
					log.Printf("Failed to get post: %v", err)
				}

				// Обработка сценария упоминания бота
				if strings.Contains(post.Message, "@"+b.user.Username) {
					var message string
					rootId := post.Id
					if post.RootId != "" {
						rootId = post.RootId
					}
					user, _, _ := b.client.GetUser(post.UserId, "")

					cmd, err := parseCommand(post.Message)
					if err != nil {
						log.Printf("Failed to parse command: %v", err)
						message = fmt.Sprintf("Привет, %s! Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужно, пиши @garroshm", user.Username)
					} else if user.Username == "casimira" {
						message = fmt.Sprintf("Сам фронтендер")
					} else {
						// TODO это временный else чтобы бот хоть что то отвечал
						message = fmt.Sprintf("Здравствуй, %s! В будущем, когда научусь, я создам таску с типом %s и заголовком ей сделаю: '%s'", user.Username, cmd.Name, cmd.Cmd)
					}

					sPost := &model.Post{
						Message:   message,
						RootId:    rootId,
						ChannelId: post.ChannelId,
					}

					err = b.SendMessage(sPost)
					if err != nil {
						log.Printf("Failed to send message: %v", err)
					}
				}
			}
		}
		fmt.Println("kill go")
	}()

	b.websocket.Listen()
}

func (b *Bot) SendMessage(post *model.Post) error {
	_, _, err := b.client.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) Close() {
	b.websocket.Close()
}
