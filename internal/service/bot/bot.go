package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"github.com/s21platform/jarvis-bot/internal/pkg/parser"
	"github.com/s21platform/jarvis-bot/internal/pkg/types"
)

type Bot struct {
	websocket      *model.WebSocketClient
	client         *model.Client4
	user           *model.User
	cmdFactory     CommandFactory
	projectMapping *config.ProjectMapping
}

func New(cfg *config.Config, cmdFactory CommandFactory) *Bot {
	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Не удалось получить информацию о пользователе: %v", err)
	}
	//log.Printf("Успешно авторизован как %s", user.Username)

	bot := &Bot{
		client:         client,
		user:           user,
		cmdFactory:     cmdFactory,
		projectMapping: config.NewProjectMapping(),
	}

	return bot
}

func (b *Bot) Listen() {
	go func() {
		for {
			if err := b.connect(); err != nil {
				log.Printf("Ошибка подключения к WebSocket: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for event := range b.websocket.EventChannel {
				if event == nil {
					log.Println("Получено пустое событие, переподключение...")
					time.Sleep(5 * time.Second)
					break
				}

				if event.EventType() == model.WebsocketEventPosted {
					post, err := parser.GetPost(event)
					if err != nil {
						log.Printf("Failed to get post: %v", err)
						continue
					}

					// Обработка сценария упоминания бота
					if strings.Contains(post.Message, "@"+b.user.Username) {
						rootId := post.Id
						if post.RootId != "" {
							rootId = post.RootId
						}
						user, _, _ := b.client.GetUser(post.UserId, "")
						channel, _, err := b.client.GetChannel(post.ChannelId, "")
						if err != nil {
							log.Printf("Failed to get channel: %v", err)
							continue
						}
						log.Println("[ CHANNEL ID ]", channel.Name)

						cmd := parser.ParseCommand(post.Message)
						var message string

						if cmd.Name == "" {
							message = fmt.Sprintf("Привет, %s! Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужна, пиши @garroshm. А весь доступный функционал ты можешь узнать по команде **help**", user.Username)
						} else {
							command := b.cmdFactory.GetCommand(cmd.Name)
							if command == nil {
								message = "Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужна, пиши @garroshm"
							} else {
								// Устанавливаем контекст команды
								cmdCtx := &types.CommandContext{
									Post:       post,
									Channel:    channel,
									User:       user,
									ProjectKey: b.projectMapping.GetProjectKey(channel.Name),
								}
								command.SetContext(cmdCtx)

								var err error
								message, err = command.Execute(cmd.Cmd)
								if err != nil {
									log.Printf("Failed to execute command: %v", err)
									message = fmt.Sprintf("Произошла ошибка при выполнении команды: %v", err)
								}
							}
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
		}
	}()
}

func (b *Bot) SendMessage(post *model.Post) error {
	_, _, err := b.client.CreatePost(post)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) connect() error {
	if b.websocket != nil {
		b.websocket.Close()
	}

	websocketURL := "wss://" + strings.TrimPrefix(b.client.URL, "https://")
	websocketClient, err := model.NewWebSocketClient4(websocketURL, b.client.AuthToken)
	if err != nil {
		return fmt.Errorf("ошибка создания WebSocket клиента: %v", err)
	}

	b.websocket = websocketClient
	b.websocket.Listen()
	return nil
}

func (b *Bot) Close() {
	if b.websocket != nil {
		b.websocket.Close()
	}
}
