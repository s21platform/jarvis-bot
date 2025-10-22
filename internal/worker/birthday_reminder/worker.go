package worker

import (
	"context"
	"log"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

type Worker struct {
	rC      RedisClient
	mC      MattermostClient
	db      DbRepo
	metrics Metrics
}

func NewWorker(rC RedisClient, db DbRepo, mC MattermostClient, metrics Metrics) *Worker {
	return &Worker{
		rC:      rC,
		mC:      mC,
		db:      db,
		metrics: metrics,
	}
}

func (w *Worker) Run() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		startTime := time.Now().UnixNano()
		w.metrics.Increment("birthday_reminder.worker.run.total")

		// Формируем ключ с текущей датой для предотвращения повторной отправки в тот же день
		today := time.Now().Format("2006-01-02")
		lockKey := "birthday_reminder_sent:" + today

		check, err := w.rC.Get(context.Background(), lockKey)
		if err != nil {
			log.Printf("Failed to get check: %v", err)
			w.metrics.Increment("birthday_reminder.redis.errors")
			continue
		}
		if check != "" {
			log.Printf("Birthday reminders already sent today (%s)", today)
			w.metrics.Increment("birthday_reminder.worker.skipped")
			continue
		}

		config, err := w.db.GetConfig(context.Background(), "birthday_reminder")
		if err != nil {
			log.Printf("Failed to get config: %v", err)
			w.metrics.Increment("birthday_reminder.db.config.errors")
			continue
		}
		birthdayReminderConfig := config.GetBirthdayReminderConfig()
		if !birthdayReminderConfig.Enabled {
			log.Println("Birthday reminder is disabled")
			w.metrics.Increment("birthday_reminder.worker.disabled")
			continue
		}
		t := time.Now()
		if t.Hour() > birthdayReminderConfig.HourEnd || t.Hour() < birthdayReminderConfig.HourStart {
			log.Println("Birthday reminder is not in the time range")
			w.metrics.Increment("birthday_reminder.worker.out_of_hours")
			continue
		}

		w.run()

		// Устанавливаем lock на 24 часа, чтобы не отправлять напоминания повторно сегодня
		err = w.rC.Set(context.Background(), lockKey, "true", 24*time.Hour)
		if err != nil {
			log.Printf("Failed to set check: %v", err)
			w.metrics.Increment("birthday_reminder.redis.errors")
			continue
		}

		w.metrics.Duration(startTime, "birthday_reminder.worker.run.duration")
	}
}

func (w *Worker) run() {
	startTime := time.Now().UnixNano()
	w.metrics.Increment("birthday_reminder.worker.process.total")

	birthdays, err := w.db.GetAllBirthdays(context.Background())
	if err != nil {
		log.Printf("Failed to get birthdays: %v", err)
		w.metrics.Increment("birthday_reminder.db.get_birthdays.errors")
		return
	}

	users, _, err := w.mC.GetUsers(0, 100, "")
	if err != nil {
		log.Printf("Failed to get users: %v", err)
		w.metrics.Increment("birthday_reminder.mattermost.get_users.errors")
		return
	}

	var processedUsers, skippedUsers, errorUsers, sentMessages int64

	// Создаем карту пользователей для быстрого доступа
	userMap := make(map[string]*model.User)
	for _, user := range users {
		if user.DeleteAt == 0 {
			userMap[user.Id] = user
		}
	}

	me, _, err := w.mC.GetMe("")
	if err != nil {
		log.Fatalf("Failed to get bot data: %v", err)
		w.metrics.Increment("birthday_reminder.mattermost.get_bot.errors")
		return
	}

	now := time.Now()
	// Обнуляем время для корректного сравнения дат
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	for userID, birthday := range birthdays {
		user, exists := userMap[userID]
		if !exists || user.DeleteAt != 0 {
			skippedUsers++
			continue
		}

		// Вычисляем следующий день рождения
		nextBirthday := time.Date(today.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.UTC)
		if nextBirthday.Before(today) {
			nextBirthday = nextBirthday.AddDate(1, 0, 0)
		}

		daysUntil := int(nextBirthday.Sub(today).Hours() / 24)

		// Проверяем, нужно ли отправлять напоминание
		if daysUntil != 7 && daysUntil != 3 && daysUntil != 0 {
			skippedUsers++
			continue
		}

		// Отправляем напоминания всем пользователям, кроме именинника
		for _, recipient := range users {
			// Пропускаем удаленных пользователей и ботов
			if recipient.DeleteAt != 0 || recipient.IsBot {
				continue
			}

			// Пропускаем именинника - он не должен получать напоминание о себе
			if recipient.Id == userID {
				continue
			}

			log.Printf("Отправляем напоминание пользователю %s о дне рождения %s (осталось дней: %d)", recipient.Username, user.Username, daysUntil)
			dmChannel, _, err := w.mC.CreateDirectChannel(me.Id, recipient.Id)
			if err != nil {
				log.Printf("Failed to create/get DM channel with %s: %v", recipient.Username, err)
				w.metrics.Increment("birthday_reminder.mattermost.create_channel.errors")
				errorUsers++
				continue
			}

			message := w.createReminderMessage(user.Username, daysUntil)
			_, _, err = w.mC.CreatePost(&model.Post{
				ChannelId: dmChannel.Id,
				Message:   message,
			})
			if err != nil {
				log.Printf("Failed to create post: %v", err)
				w.metrics.Increment("birthday_reminder.mattermost.create_post.errors")
				errorUsers++
				continue
			}

			sentMessages++
			processedUsers++
		}
	}

	w.metrics.Count("birthday_reminder.users.processed", processedUsers)
	w.metrics.Count("birthday_reminder.users.skipped", skippedUsers)
	w.metrics.Count("birthday_reminder.users.errors", errorUsers)
	w.metrics.Count("birthday_reminder.messages.sent", sentMessages)
	w.metrics.Duration(startTime, "birthday_reminder.worker.process.duration")
}

func (w *Worker) createReminderMessage(username string, daysUntil int) string {
	switch daysUntil {
	case 0:
		return "🎉 Сегодня день рождения у @" + username + "!"
	case 3:
		return "🎈 Через 3 дня день рождения у @" + username + "! Готовь поздравления! Если хочешь присоединиться, у нас есть сбор средств на поздравление!\nСчет единый на все дни рождения. Если хочешь поучаствовать в сборе, оставь, пожалуйста, комментарий в переводе, на чей день рождения отправляешь ✨\nhttps://www.tbank.ru/cf/AvF05i66rx8"
	case 7:
		return "📅 Через неделю день рождения у @" + username + "!"
	default:
		return ""
	}
}
