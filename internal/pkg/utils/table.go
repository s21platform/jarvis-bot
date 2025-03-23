package utils

import "strings"

// cleanMarkdown удаляет MD разметку заголовков из текста и обрабатывает переносы строк
func cleanMarkdown(text string) string {
	// Заменяем переносы строк на пробелы
	text = strings.ReplaceAll(text, "\n", " ")

	// Удаляем символы # в начале строки
	text = strings.TrimLeft(text, "# ")

	// Удаляем множественные пробелы
	text = strings.Join(strings.Fields(text), " ")

	return text
}

// cleanCell подготавливает ячейку для отображения в таблице
func cleanCell(text string) string {
	// Очищаем MD разметку
	text = cleanMarkdown(text)

	// Если в тексте есть "|", заменяем на "-"
	text = strings.ReplaceAll(text, "|", "-")

	// Если в тексте есть ссылка "читать далее", оставляем её в конце
	if idx := strings.Index(text, "[читать далее]"); idx != -1 {
		linkStart := idx
		linkEnd := strings.Index(text[linkStart:], ")")
		if linkEnd != -1 {
			linkEnd += linkStart + 1 // +1 чтобы включить закрывающую скобку
			beforeLink := strings.TrimSpace(text[:linkStart])
			link := text[linkStart:linkEnd]
			return beforeLink + " " + link
		}
	}

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
			// Подготавливаем ячейку для отображения
			cleanedCell := cleanCell(cell)
			builder.WriteString(" " + cleanedCell + " |")
		}
		builder.WriteString("\n")
	}

	return builder.String()
}
