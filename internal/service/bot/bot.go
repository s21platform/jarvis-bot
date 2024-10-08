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
	dbR       DbRepo
}

func New(cfg *config.Config, dbR DbRepo) *Bot {
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
		dbR:       dbR,
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
					channel, _, err := b.client.GetChannel(post.ChannelId, "")
					if err != nil {
						log.Printf("Failed to get channel: %v", err)
						continue
					}
					log.Println("[ CHANNEL ID ]", channel.Name)

					cmd := parseCommand(post.Message)

					switch cmd.Name {
					case "help":
						message = "Привет! Это мой --help. Сейчас я знаю команды:\n" +
							"- feature <name_of_feature>\t[Создает фичу для текущего сервиса]\n" +
							"- bug <name_of_bug>\t[Создает баг для текущего сервиса]\n" +
							"- my \t[Показывет таски инициатора в данном сервисе]\n"
					case "feature":
						id, err := b.dbR.CreateTask(channel.Name, "feature", cmd.Cmd, post.UserId)
						if err != nil {
							log.Printf("Failed to create feature: %v", err)
							continue
						}
						message = fmt.Sprintf("Создал feature с заголовком: %s, для сервиса %s с id: %d", cmd.Cmd, channel.Name, id)
					case "bug":
						message = fmt.Sprintf("В будущем, когда научусь, я создам таску с типом **bug** и заголовком ей сделаю: '%s'", cmd.Cmd)
					case "my":
						tasks, err := b.dbR.GetTasksByUUID(post.UserId, channel.Name)
						if err != nil {
							log.Printf("Failed to get tasks: %v", err)
							continue
						}
						message = CreateTable([]string{"ID", "Таска", "Описание", "Тип"}, convertModelToString(tasks))
					case "tasks":
						tasks, err := b.dbR.GetTasksByChannel(channel.Name)
						if err != nil {
							log.Printf("Failed to get tasks: %v", err)
							continue
						}
						message = CreateTable([]string{"ID", "Исполнитель", "Таска", "Описание", "Тип"}, convertModelAllTasksToString(tasks))
					default:
						message = fmt.Sprintf("Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужна, пиши @garroshm")
					}

					if cmd.Name == "" {
						message = fmt.Sprintf("Привет, %s! Такая команда мне еще не знакома. Если ты считаешь, что такая команда нужно, пиши @garroshm. А весь доступный функционал ты можешь узнать по команде **help**", user.Username)
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
