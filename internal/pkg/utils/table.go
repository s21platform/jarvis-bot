package utils

import "strings"

// cleanMarkdown удаляет MD разметку заголовков из текста
func cleanMarkdown(text string) string {
	// Удаляем символы # в начале строки
	text = strings.TrimLeft(text, "# ")
	return text
}

// CreateTable создает строковое представление таблицы
func CreateTable(headers []string, rows [][]string) string {
	var builder strings.Builder

	// Создаем строку заголовков
	builder.WriteString("|")
	for _, header := range headers {
		builder.WriteString(" " + header + " |")
	}
	builder.WriteString("\n")

	// Создаем разделитель
	builder.WriteString("|")
	for range headers {
		builder.WriteString("------------|")
	}
	builder.WriteString("\n")

	// Добавляем строки данных
	for _, row := range rows {
		builder.WriteString("|")
		for _, cell := range row {
			// Очищаем MD разметку перед добавлением в таблицу
			cleanCell := cleanMarkdown(cell)
			builder.WriteString(" " + cleanCell + " |")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
