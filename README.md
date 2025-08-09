# Jarvis Bot

Бот для Mattermost с модульной архитектурой и поддержкой различных команд.

## Структура проекта

```
internal/
├── command/              # Команды бота
│   ├── cred/            # Команда для работы с учетными данными
│   ├── help/            # Команда помощи
│   ├── news/            # Команда для работы с новостями
│   └── factory.go       # Фабрика команд
├── config/              # Конфигурация
├── model/               # Модели данных
├── pkg/                 # Общие пакеты
│   ├── parser/         # Парсер команд
│   ├── types/          # Общие типы и интерфейсы
│   └── utils/          # Утилиты
├── repository/          # Слой работы с БД
│   └── postgres/       # Реализация для PostgreSQL
└── service/            # Сервисный слой
    └── bot/            # Основная логика бота
        └── contract.go # Интерфейсы бота
```

## Как добавить новую команду

1. Создайте новый пакет в директории `internal/command/`:
```bash
mkdir internal/command/mycommand
```

2. Создайте файл `command.go` в новом пакете со следующей структурой:
```go
package mycommand

import (
    "github.com/s21platform/jarvis-bot/internal/pkg/types"
    "github.com/s21platform/jarvis-bot/internal/service/bot"
)

type Command struct {
    types.BaseCommand
    db bot.DbRepo // если нужна работа с БД
    // другие необходимые поля
}

func NewCommand(db bot.DbRepo) *Command {
    return &Command{
        BaseCommand: types.NewBaseCommand("mycommand", "Описание команды"),
        db:         db,
    }
}

// Реализация методов интерфейса Command
func (c *Command) Execute(args string) (string, error) {
    // Логика выполнения команды
    return "Результат выполнения", nil
}
```

3. Добавьте команду в фабрику (`internal/command/factory.go`):
```go
import (
    "github.com/s21platform/jarvis-bot/internal/command/mycommand"
)

func NewFactory(client *model.Client4, user *model.User, db bot.DbRepo) *Factory {
    // ...
    myCmd := mycommand.NewCommand(db)
    commands := []types.Command{
        // ...
        myCmd,
    }
    // ...
}
```

## Доступные команды

### help
Показывает список доступных команд и их описание.
```
@jarvis help
```

### cred
Получение учетных данных для сервиса.
```
@jarvis cred <service_name>
```

### news
Управление новостями.
```
@jarvis news add      # Добавить новость (в треде)
@jarvis news delete   # Удалить новость (в треде)
@jarvis news show     # Показать новости за последние 3 недели
```

Особенности команды news:
- При добавлении новости бот ставит реакцию ⚡ на сообщение
- При удалении новости реакция автоматически убирается
- Длинные новости (>80 символов) сокращаются с добавлением ссылки "читать далее"

### create
Создание задач и багов в Jira.
```
@jarvis create task <заголовок>  # Создать задачу 💡
@jarvis create bug <заголовок>   # Создать баг 🐞
```

Особенности команды create:
- Автоматически определяет проект Jira на основе канала
- Возвращает кликабельную ссылку на созданную задачу/баг

#### Настройка маппинга каналов в проекты Jira

Для работы команды `create` необходимо настроить соответствие между каналами Mattermost и проектами Jira. Это делается в файле `internal/config/project_mapping.go`:

```go
func NewProjectMapping() *ProjectMapping {
    return &ProjectMapping{
        ChannelToProject: map[string]string{
            // Примеры маппинга:
            "chat-team-backend": "BACK",     // Канал -> Ключ проекта в Jira
            "chat-team-frontend": "FRONT",
            
            // Community
            model.COMMUNITY_SERVICE_CHANNEL:        model.COMMUNITY_PROJECT,
            model.COMMUNITY_SERVICE_NEWS_CHANNEL:   model.COMMUNITY_PROJECT,
            model.COMMUNITY_SERVICE_PUBLIC_CHANNEL: model.COMMUNITY_PROJECT,
            
            // Добавьте свои маппинги здесь
        },
    }
}

## Разработка

### Требования
- Go 1.21+
- PostgreSQL 14+
- Mattermost Server
- Jira Server/Cloud

### Конфигурация
Создайте файл `.env` на основе `.env.example`:
```bash
cp .env.example .env
```

### Запуск
```bash
go run cmd/bot/main.go
```

## Архитектура

Проект следует принципам чистой архитектуры:
1. Зависимости направлены внутрь
2. Бизнес-логика не зависит от деталей реализации
3. Интерфейсы определены в слое их использования

### Основные интерфейсы

#### Command (types.Command)
```go
type Command interface {
    Execute(args string) (string, error)
    Name() string
    Description() string
    SetContext(ctx *CommandContext)
}
```

#### DbRepo (bot.DbRepo)
```go
type DbRepo interface {
    GetCred(ctx context.Context, name string) (model.Data, error)
    AddNews(ctx context.Context, threadID, channel, messageID string) error
    DeleteNews(ctx context.Context, threadID string) error
    GetNews(ctx context.Context) ([]model.News, error)
    IsNewsExists(ctx context.Context, threadID string) (bool, error)
    Close()
}
```
