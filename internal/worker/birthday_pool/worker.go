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
	metrics     Metrics
}

func NewWorker(rC RedisClient, db DbRepo, mC MattermostClient, cfg *config.Config, metrics Metrics) *Worker {
	return &Worker{
		rC:          rC,
		mC:          mC,
		db:          db,
		callbackUrl: cfg.Service.Url,
		metrics:     metrics,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		startTime := time.Now().UnixNano()
		w.metrics.Increment("birthday_pool.worker.run.total")

		check, err := w.rC.Get(context.Background(), "check")
		if err != nil {
			log.Printf("Failed to get check: %v", err)
			w.metrics.Increment("birthday_pool.redis.errors")
			continue
		}
		if check != "" {
			w.metrics.Increment("birthday_pool.worker.skipped")
			continue
		}

		config, err := w.db.GetConfig(context.Background(), "birthday_pool")
		if err != nil {
			log.Printf("Failed to get config: %v", err)
			w.metrics.Increment("birthday_pool.db.config.errors")
			return
		}
		birthdayPoolConfig := config.GetBirthdayPoolConfig()
		if !birthdayPoolConfig.Enabled {
			log.Println("Birthday pool is disabled")
			w.metrics.Increment("birthday_pool.worker.disabled")
			return
		}

		w.run()

		err = w.rC.Set(context.Background(), "check", "true", time.Duration(birthdayPoolConfig.PendingIntervalMinutes)*time.Minute)
		if err != nil {
			log.Printf("Failed to set check: %v", err)
			w.metrics.Increment("birthday_pool.redis.errors")
			continue
		}

		w.metrics.Duration(startTime, "birthday_pool.worker.run.duration")
	}
}

func (w *Worker) run() {
	startTime := time.Now().UnixNano()
	w.metrics.Increment("birthday_pool.worker.process.total")

	config, err := w.db.GetConfig(context.Background(), "birthday_pool")
	if err != nil {
		log.Printf("Failed to get config: %v", err)
		w.metrics.Increment("birthday_pool.db.config.errors")
		return
	}
	birthdayPoolConfig := config.GetBirthdayPoolConfig()
	if !birthdayPoolConfig.Enabled {
		log.Println("Birthday pool is disabled")
		w.metrics.Increment("birthday_pool.worker.disabled")
		return
	}

	t := time.Now()
	if t.Hour() > birthdayPoolConfig.HourEnd || t.Hour() < birthdayPoolConfig.HourStart {
		log.Println("Birthday pool is not in the time range")
		w.metrics.Increment("birthday_pool.worker.out_of_hours")
		return
	}

	users, _, err := w.mC.GetUsers(0, 100, "")
	if err != nil {
		w.metrics.Increment("birthday_pool.mattermost.get_users.errors")
		panic(err)
	}
	w.metrics.Count("birthday_pool.users.total", int64(len(users)))

	var processedUsers, skippedUsers, errorUsers, sentMessages int64

	for _, user := range users {
		if user.DeleteAt != 0 {
			skippedUsers++
			continue
		}

		birthday, err := w.db.GetBirthday(context.Background(), user.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Failed to get birthday: %v", err)
			w.metrics.Increment("birthday_pool.db.get_birthday.errors")
			errorUsers++
			continue
		}

		if birthday != nil {
			log.Println("Birthday: ", user.Username, birthday)
			skippedUsers++
			continue
		}

		if len(birthdayPoolConfig.WhiteList) > 0 {
			if !slices.Contains(birthdayPoolConfig.WhiteList, user.Username) {
				log.Println("User is not in the white list")
				skippedUsers++
				continue
			}
		}

		me, _, err := w.mC.GetMe("")
		if err != nil {
			log.Fatalf("Не удалось получить данные о боте: %v", err)
			w.metrics.Increment("birthday_pool.mattermost.get_bot.errors")
			errorUsers++
			continue
		}

		dmChannel, _, err := w.mC.CreateDirectChannel(me.Id, user.Id)
		if err != nil {
			log.Printf("Не удалось создать/получить DM-канал с %s: %v", user.Username, err)
			w.metrics.Increment("birthday_pool.mattermost.create_channel.errors")
			errorUsers++
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
			w.metrics.Increment("birthday_pool.mattermost.create_post.errors")
			errorUsers++
			continue
		}

		sentMessages++
		processedUsers++
	}

	// Отправляем итоговые метрики
	w.metrics.Count("birthday_pool.users.processed", processedUsers)
	w.metrics.Count("birthday_pool.users.skipped", skippedUsers)
	w.metrics.Count("birthday_pool.users.errors", errorUsers)
	w.metrics.Count("birthday_pool.messages.sent", sentMessages)
	w.metrics.Duration(startTime, "birthday_pool.worker.process.duration")
}
