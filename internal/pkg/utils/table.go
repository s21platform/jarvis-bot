package utils

import "strings"

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
			builder.WriteString(" " + cell + " |")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
