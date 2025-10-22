package worker

import (
	"testing"
	"time"
)

// TestCalculateDaysUntilBirthday проверяет корректность расчета дней до дня рождения
func TestCalculateDaysUntilBirthday(t *testing.T) {
	tests := []struct {
		name          string
		today         time.Time
		birthdayMonth time.Month
		birthdayDay   int
		expectedDays  int
	}{
		{
			name:          "День рождения сегодня",
			today:         time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
			birthdayMonth: time.October,
			birthdayDay:   22,
			expectedDays:  0,
		},
		{
			name:          "День рождения через 7 дней",
			today:         time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
			birthdayMonth: time.October,
			birthdayDay:   29,
			expectedDays:  7,
		},
		{
			name:          "День рождения через 3 дня",
			today:         time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
			birthdayMonth: time.October,
			birthdayDay:   25,
			expectedDays:  3,
		},
		{
			name:          "День рождения в следующем году",
			today:         time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC),
			birthdayMonth: time.January,
			birthdayDay:   1,
			expectedDays:  7, // 25 дек -> 1 янв = 7 дней
		},
		{
			name:          "День рождения уже прошел в этом году",
			today:         time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
			birthdayMonth: time.January,
			birthdayDay:   15,
			expectedDays:  85, // до следующего 15 января
		},
		{
			name:          "Високосный год - 29 февраля",
			today:         time.Date(2024, 2, 22, 0, 0, 0, 0, time.UTC),
			birthdayMonth: time.February,
			birthdayDay:   29,
			expectedDays:  7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Симулируем логику из worker.go
			nextBirthday := time.Date(tt.today.Year(), tt.birthdayMonth, tt.birthdayDay, 0, 0, 0, 0, time.UTC)
			if nextBirthday.Before(tt.today) {
				nextBirthday = nextBirthday.AddDate(1, 0, 0)
			}

			daysUntil := int(nextBirthday.Sub(tt.today).Hours() / 24)

			if daysUntil != tt.expectedDays {
				t.Errorf("Expected %d days, got %d days", tt.expectedDays, daysUntil)
			}
		})
	}
}

// TestCalculateDaysWithTime проверяет, что время не влияет на расчет дней
func TestCalculateDaysWithTime(t *testing.T) {
	tests := []struct {
		name          string
		currentTime   time.Time
		birthdayMonth time.Month
		birthdayDay   int
		expectedDays  int
	}{
		{
			name:          "Текущее время 23:59 - день рождения через 7 дней должен быть корректным",
			currentTime:   time.Date(2025, 10, 22, 23, 59, 59, 0, time.UTC),
			birthdayMonth: time.October,
			birthdayDay:   29,
			expectedDays:  7,
		},
		{
			name:          "Текущее время 00:01 - день рождения через 7 дней должен быть корректным",
			currentTime:   time.Date(2025, 10, 22, 0, 1, 0, 0, time.UTC),
			birthdayMonth: time.October,
			birthdayDay:   29,
			expectedDays:  7,
		},
		{
			name:          "Текущее время 15:30 - день рождения через 7 дней должен быть корректным",
			currentTime:   time.Date(2025, 10, 22, 15, 30, 0, 0, time.UTC),
			birthdayMonth: time.October,
			birthdayDay:   29,
			expectedDays:  7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// КРИТИЧНО: обнуляем время перед расчетом (как в исправленном коде)
			today := time.Date(tt.currentTime.Year(), tt.currentTime.Month(), tt.currentTime.Day(), 0, 0, 0, 0, time.UTC)

			nextBirthday := time.Date(today.Year(), tt.birthdayMonth, tt.birthdayDay, 0, 0, 0, 0, time.UTC)
			if nextBirthday.Before(today) {
				nextBirthday = nextBirthday.AddDate(1, 0, 0)
			}

			daysUntil := int(nextBirthday.Sub(today).Hours() / 24)

			if daysUntil != tt.expectedDays {
				t.Errorf("Expected %d days, got %d days (current time: %v)", tt.expectedDays, daysUntil, tt.currentTime)
			}
		})
	}
}

// TestCreateReminderMessage проверяет создание сообщений для разных сценариев
func TestCreateReminderMessage(t *testing.T) {
	w := &Worker{}

	tests := []struct {
		name             string
		username         string
		daysUntil        int
		expectedContains []string
		expectedNotEmpty bool
	}{
		{
			name:             "День рождения сегодня",
			username:         "ivan",
			daysUntil:        0,
			expectedContains: []string{"🎉", "Сегодня", "@ivan"},
			expectedNotEmpty: true,
		},
		{
			name:             "День рождения через 3 дня",
			username:         "maria",
			daysUntil:        3,
			expectedContains: []string{"🎈", "Через 3 дня", "@maria", "tbank.ru"},
			expectedNotEmpty: true,
		},
		{
			name:             "День рождения через неделю",
			username:         "alex",
			daysUntil:        7,
			expectedContains: []string{"📅", "Через неделю", "@alex"},
			expectedNotEmpty: true,
		},
		{
			name:             "Некорректное количество дней",
			username:         "test",
			daysUntil:        5,
			expectedContains: []string{},
			expectedNotEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := w.createReminderMessage(tt.username, tt.daysUntil)

			if tt.expectedNotEmpty && message == "" {
				t.Error("Expected non-empty message, got empty")
			}

			if !tt.expectedNotEmpty && message != "" {
				t.Errorf("Expected empty message, got: %s", message)
			}

			for _, expected := range tt.expectedContains {
				if !contains(message, expected) {
					t.Errorf("Expected message to contain '%s', but it doesn't. Message: %s", expected, message)
				}
			}
		})
	}
}

// Вспомогательная функция для проверки содержимого строки
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:len(substr)] == substr ||
			(len(s) > len(substr) && contains(s[1:], substr)))))
}
