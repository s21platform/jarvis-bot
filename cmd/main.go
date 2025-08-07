package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
	"log"
)

func main() {
	cfg := config.MustLoadConfig()
	log.Println(cfg.Url)
	client := model.NewAPIv4Client(cfg.Url)
	client.SetOAuthToken(cfg.Token)

	users, _, err := client.GetUsers(0, 100, "")
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		if user.DeleteAt != 0 {
			log.Println("User deleted: ", user.Username)
		} else {
			log.Println(user)
		}

		if user.Username == "garroshm" || user.Username == "violasab" || user.Username == "whiteyes" {
			me, _, err := client.GetMe("")
			if err != nil {
				log.Fatalf("Не удалось получить данные о боте: %v", err)
			}
			dmChannel, _, err := client.CreateDirectChannel(me.Id, user.Id)
			if err != nil {
				log.Printf("Не удалось создать/получить DM-канал с %s: %v", user.Username, err)
				continue
			}
			log.Println("here")
			_, _, err = client.CreatePost(&model.Post{
				ChannelId: dmChannel.Id,
				Message:   "Привет! 👋 Мы в Space-21 собираем даты дней рождения команды, чтобы поздравлять друг друга вовремя 🎉\nЗаполни, пожалуйста, короткую форму — это займёт минуту 🙌",
				Props: map[string]interface{}{
					"attachments": []map[string]interface{}{
						{
							"actions": []map[string]interface{}{
								{
									"name": "Открыть форму",
									"integration": map[string]interface{}{
										"url": "https://12c3672ee943.ngrok-free.app/show-birthday-dialog-window",
										"context": map[string]interface{}{
											"action": "open_modal_birthday",
										},
									},
									"type": "button",
								},
							},
						},
					},
				},
			})
			if err != nil {
				log.Printf("Failed to create post: %v", err)
			}
		}

	}

}
