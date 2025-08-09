package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"

	"github.com/s21platform/jarvis-bot/internal/config"
	api "github.com/s21platform/jarvis-bot/internal/generated"
	internal_model "github.com/s21platform/jarvis-bot/internal/model"
)

type Handler struct {
	bot         *model.User
	client      *model.Client4
	dbR         DatabaseRepo
	callbackUrl string
}

func New(cfg *config.Config, dbR DatabaseRepo, callbackUrl string) *Handler {
	client := model.NewAPIv4Client(cfg.Bot.Url)
	client.SetOAuthToken(cfg.Bot.Token)

	bot, _, err := client.GetMe("")
	if err != nil {
		log.Fatalf("Не удалось получить информацию о пользователе: %v", err)
	}
	return &Handler{
		bot:         bot,
		client:      client,
		dbR:         dbR,
		callbackUrl: callbackUrl,
	}
}

func (h *Handler) PostSaveBirthday(w http.ResponseWriter, r *http.Request) {
	var in api.HandleSavingBirthday
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, _, err := h.client.GetUser(in.UserId, "")
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		h.sendCallbackMessage(in.UserId, "Произошла ошибка сохранения данных. Напиши пожалуйста @garroshm об этом")
		return
	}

	dayInt, err := strconv.ParseInt(in.Submission.Day, 10, 64)
	if err != nil {
		log.Println(err.Error())
		h.sendCallbackMessage(in.UserId, "Произошла ошибка сохранения данных. Напиши пожалуйста @garroshm об этом")
		return
	}

	monthInt, err := strconv.ParseInt(in.Submission.Month, 10, 64)
	if err != nil {
		log.Println(err.Error())
		h.sendCallbackMessage(in.UserId, "Произошла ошибка сохранения данных. Напиши пожалуйста @garroshm об этом")
		return
	}

	var yearPtr *int64
	if in.Submission.Year != nil {
		year := *in.Submission.Year
		yearInt, err := strconv.ParseInt(year, 10, 64)
		if err != nil {
			log.Println(err.Error())
			yearInt = 0
		}
		if yearInt < 1900 || yearInt >= int64(time.Now().Year()) {
			log.Println("illegal year")
			yearInt = 0
		}
		if yearInt == 0 {
			h.sendCallbackMessage(in.UserId, "Произошла ошибка сохранения данных. Напиши пожалуйста @garroshm об этом")
			yearPtr = nil
		} else {
			yearPtr = &yearInt
		}
	}

	var channelId string
	if in.ChannelId != nil {
		channelId = *in.ChannelId
	}

	err = h.dbR.SetBirthday(r.Context(), &internal_model.Birthday{
		Day:   dayInt,
		Month: monthInt,
		Year:  yearPtr,
	}, in.UserId, channelId, user.Username)

	if err != nil {
		log.Println(err.Error())
		h.sendCallbackMessage(in.UserId, "Произошла ошибка сохранения данных. Напиши пожалуйста @garroshm об этом")
		return
	}

	_, _, err = h.client.CreatePost(&model.Post{
		ChannelId: channelId,
		UserId:    in.UserId,
		Message:   "Спасибо! 🎉 Записали твою дату рождения — теперь точно не забудем поздравить! 🥳",
	})
	if err != nil {
		log.Println(err.Error())
		h.sendCallbackMessage(in.UserId, "Произошла ошибка сохранения данных. Напиши пожалуйста @garroshm об этом")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) sendCallbackMessage(userId string, message string) {
	_, _, err := h.client.CreatePost(&model.Post{
		Message: message,
		UserId:  userId,
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func (h *Handler) PostShowBirthdayDialogWindow(w http.ResponseWriter, r *http.Request) {
	log.Println("PostShowBirthdayDialogWindow")
	var in api.OpenBirthdayWindow
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("userId: ", in.UserId)

	birthday, err := h.dbR.GetBirthday(r.Context(), in.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err.Error())
		return
	}
	if birthday != nil {
		// TODO Отправить сообщение, что дата уже установлена
		log.Println("birthday: ", birthday)
		log.Println("already have birthday")
		return
	}

	_, err = h.client.OpenInteractiveDialog(model.OpenDialogRequest{
		TriggerId: in.TriggerId,
		URL:       h.callbackUrl + "/save-birthday",
		Dialog: model.Dialog{
			Title: "Твой день рождения",
			Elements: []model.DialogElement{
				{
					Name:        "day",
					DisplayName: "День рождения",
					Type:        "select",
					Options:     generateDays(),
				},
				{
					Name:        "month",
					DisplayName: "Месяц рождения",
					Type:        "select",
					Options:     generateMonths(),
				},
				{
					Name:        "year",
					DisplayName: "Год Рождения",
					Type:        "select",
					Optional:    true,
					Default:     "0",
					HelpText:    "Вводить год не обязательно, при любом случае мы не будем нотифицировать о количестве лет",
					Options:     generateYears(),
				},
			},
		},
	})
	if err != nil {
		log.Println(err.Error())
	}
}

func generateDays() []*model.PostActionOptions {
	var res []*model.PostActionOptions
	for i := range 31 {
		res = append(res, &model.PostActionOptions{
			Text:  fmt.Sprintf("%02d", i+1),
			Value: fmt.Sprintf("%d", i+1),
		})
	}
	return res
}

func generateMonths() []*model.PostActionOptions {
	var res []*model.PostActionOptions
	for i, m := range []string{
		"Январь",
		"Февраль",
		"Март",
		"Апрель",
		"Май",
		"Июнь",
		"Июль",
		"Август",
		"Сентябрь",
		"Октябрь",
		"Ноябрь",
		"Декабрь",
	} {
		res = append(res, &model.PostActionOptions{
			Text:  m,
			Value: fmt.Sprintf("%d", i+1),
		})
	}
	return res
}

func generateYears() []*model.PostActionOptions {
	var res []*model.PostActionOptions
	yearStart := 1900
	yearFinish := time.Now().Add(-1 * 12 * 365 * 24 * time.Hour).Year()
	for i := yearFinish; i >= yearStart; i-- {
		res = append(res, &model.PostActionOptions{
			Text:  fmt.Sprintf("%02d", i),
			Value: fmt.Sprintf("%d", i),
		})
	}
	return res
}
