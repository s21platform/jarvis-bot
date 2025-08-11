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
	metrics        Metrics
}

func New(cfg *config.Config, cmdFactory CommandFactory, metrics Metrics) *Bot {
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
		metrics:        metrics,
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
					b.metrics.Increment("bot.websocket.events.empty")
					time.Sleep(5 * time.Second)
					break
				}

				b.metrics.Increment("bot.websocket.events.total")

				if event.EventType() == model.WebsocketEventPosted {
					b.metrics.Increment("bot.websocket.events.posted")
					post, err := parser.GetPost(event)
					if err != nil {
						log.Printf("Failed to get post: %v", err)
						b.metrics.Increment("bot.errors.parse_post")
						continue
					}

					// Обработка сценария упоминания бота
					if strings.Contains(post.Message, "@"+b.user.Username) {
						b.metrics.Increment("bot.messages.mentions")
						rootId := post.Id
						if post.RootId != "" {
							rootId = post.RootId
						}
						user, _, _ := b.client.GetUser(post.UserId, "")
						channel, _, err := b.client.GetChannel(post.ChannelId, "")
						if err != nil {
							log.Printf("Failed to get channel: %v", err)
							b.metrics.Increment("bot.errors.get_channel")
							continue
						}
						log.Println("[ CHANNEL ID ]", channel.Name)

						startTime := time.Now().UnixNano()
						b.metrics.Increment("bot.commands.total")

						cmd := parser.ParseCommand(post.Message)
						var message string

						if cmd.Name == "" {
							b.metrics.Increment("bot.commands.unknown")
							message = fmt.Sprintf("Привет, %s! Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужна, пиши @garroshm. А весь доступный функционал ты можешь узнать по команде **help**", user.Username)
						} else {
							command := b.cmdFactory.GetCommand(cmd.Name)
							if command == nil {
								b.metrics.Increment("bot.commands.unknown")
								message = "Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужна, пиши @garroshm"
							} else {
								b.metrics.Increment("bot.commands." + cmd.Name)
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
									b.metrics.Increment("bot.commands." + cmd.Name + ".errors")
									message = fmt.Sprintf("Произошла ошибка при выполнении команды: %v", err)
								}
								b.metrics.Duration(startTime, "bot.commands."+cmd.Name+".duration")
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
							b.metrics.Increment("bot.errors.send_message")
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
	startTime := time.Now().UnixNano()
	b.metrics.Increment("bot.websocket.connect.total")

	if b.websocket != nil {
		b.websocket.Close()
	}

	websocketURL := "wss://" + strings.TrimPrefix(b.client.URL, "https://")
	websocketClient, err := model.NewWebSocketClient4(websocketURL, b.client.AuthToken)
	if err != nil {
		b.metrics.Increment("bot.websocket.connect.errors")
		return fmt.Errorf("ошибка создания WebSocket клиента: %v", err)
	}

	b.websocket = websocketClient
	b.websocket.Listen()

	b.metrics.Duration(startTime, "bot.websocket.connect.duration")
	b.metrics.Increment("bot.websocket.connect.success")
	return nil
}

func (b *Bot) Close() {
	if b.websocket != nil {
		b.websocket.Close()
	}
}
