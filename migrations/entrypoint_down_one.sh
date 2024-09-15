#!/bin/bash

# Применим миграции
# ВАЖНО! Для локального запуска вызвать из корня проекта, предварительно подгрузив переменные окружения в терминал командой  `set -a; source <путь к .env файлу>; set +a`
# Установка goose: go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir ./migrations postgres "user=$JARVIS_BOT_POSTGRES_USER password=$JARVIS_BOT_POSTGRES_PASSWORD dbname=$JARVIS_BOT_POSTGRES_DB host=$JARVIS_BOT_POSTGRES_HOST port=$JARVIS_BOT_POSTGRES_PORT sslmode=disable" down