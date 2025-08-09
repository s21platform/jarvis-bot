package worker

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"slices"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/s21platform/jarvis-bot/internal/config"
)

type Worker struct {
	rC          RedisClient
	mC          MattermostClient
	db          DbRepo
	callbackUrl string
}

func NewWorker(rC RedisClient, db DbRepo, mC MattermostClient, cfg *config.Config) *Worker {
	return &Worker{
		rC:          rC,
		mC:          mC,
		db:          db,
		callbackUrl: cfg.Service.Url,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		check, err := w.rC.Get(context.Background(), "check")
		if err != nil {
			log.Printf("Failed to get check: %v", err)
			continue
		}
		if check != "" {
			continue
		}

		config, err := w.db.GetConfig(context.Background(), "birthday_pool")
		if err != nil {
			log.Printf("Failed to get config: %v", err)
			return
		}
		birthdayPoolConfig := config.GetBirthdayPoolConfig()
		if !birthdayPoolConfig.Enabled {
			log.Println("Birthday pool is disabled")
			return
		}

		w.run()

		err = w.rC.Set(context.Background(), "check", "true", time.Duration(birthdayPoolConfig.PendingIntervalMinutes)*time.Minute)
		if err != nil {
			log.Printf("Failed to set check: %v", err)
			continue
		}
	}
}

func (w *Worker) run() {
	config, err := w.db.GetConfig(context.Background(), "birthday_pool")
	if err != nil {
		log.Printf("Failed to get config: %v", err)
		return
	}
	birthdayPoolConfig := config.GetBirthdayPoolConfig()
	if !birthdayPoolConfig.Enabled {
		log.Println("Birthday pool is disabled")
		return
	}

	t := time.Now()
	if t.Hour() > birthdayPoolConfig.HourEnd || t.Hour() < birthdayPoolConfig.HourStart {
		log.Println("Birthday pool is not in the time range")
		return
	}
	users, _, err := w.mC.GetUsers(0, 100, "")
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		if user.DeleteAt != 0 {
			continue
		}

		birthday, err := w.db.GetBirthday(context.Background(), user.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to get birthday: %v", err)
			continue
		}

		if birthday != nil {
			log.Println("Birthday: ", user.Username, birthday)
			continue
		}

		if len(birthdayPoolConfig.WhiteList) > 0 {
			if !slices.Contains(birthdayPoolConfig.WhiteList, user.Username) {
				log.Println("User is not in the white list")
				continue
			}
		}
		me, _, err := w.mC.GetMe("")
		if err != nil {
			log.Fatalf("Не удалось получить данные о боте: %v", err)
		}
		dmChannel, _, err := w.mC.CreateDirectChannel(me.Id, user.Id)
		if err != nil {
			log.Printf("Не удалось создать/получить DM-канал с %s: %v", user.Username, err)
			continue
		}
		log.Println("Created DM channel: ", dmChannel.Id, user.Username, user.Id)
		_, _, err = w.mC.CreatePost(&model.Post{
			ChannelId: dmChannel.Id,
			Message:   "Привет! 👋 Мы в Space-21 собираем даты дней рождения команды, чтобы поздравлять друг друга вовремя 🎉\nЗаполни, пожалуйста, короткую форму — это займёт минуту 🙌",
			Props: map[string]interface{}{
				"attachments": []map[string]interface{}{
					{
						"actions": []map[string]interface{}{
							{
								"name": "Открыть форму",
								"integration": map[string]interface{}{
									"url": w.callbackUrl + "/show-birthday-dialog-window",
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
