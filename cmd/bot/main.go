package main

import (
	"encoding/json"
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"log"
	"strings"
	"time"
)

func main() {
	cfg := config.MustLoadConfig()
	client := model.NewAPIv4Client(cfg.Url)
	client.SetOAuthToken(cfg.Token)

	user, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Не удалось получить информацию о пользователе: %v", err)
	}
	log.Printf("Успешно авторизован как %s", user.Username)

	websocketClient, err := model.NewWebSocketClient4("wss://"+cfg.Url[len("https://"):], client.AuthToken)
	if err != nil {
		log.Fatalf("Ошибка подключения к WebSocket: %v", err)
	}
	defer websocketClient.Close()

	go func() {
		time.Sleep(2 * time.Second)
		log.Println("in gor", websocketClient.EventChannel)
		for event := range websocketClient.EventChannel {
			fmt.Println("in loop")
			if event.EventType() == model.WebsocketEventPosted {
				rawPost := event.GetData()["post"].(string)
				post := model.Post{}
				err := json.Unmarshal([]byte(rawPost), &post)
				if err != nil {
					log.Printf("Ошибка десериализации поста: %v", err)
					continue
				}
				message := post.Message
				channelID := post.ChannelId

				// Проверяем, упомянут ли бот в сообщении
				if strings.Contains(message, "@"+user.Username) {
					responseMessage := fmt.Sprintf("Привет! Пи-р garroshm создал меня видимо просто поржать, так как я сейчас еще нчего не умею. Но спасибо что нашли меня, спасите!")
					responsePost := &model.Post{
						ChannelId: channelID,
						Message:   responseMessage,
					}
					_, _, err := client.CreatePost(responsePost)
					if err != nil {
						log.Printf("Ошибка при отправке сообщения: %v", err)
					}
				}

			}

		}
	}()

	websocketClient.Listen()
	select {}
}
