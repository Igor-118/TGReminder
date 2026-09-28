# TGReminder

Telegram-бот для напоминалок: добавляешь напоминание — бот присылает его в нужное время.

## Стек

- Go 1.26
- [telegram-bot-api v5](https://github.com/go-telegram-bot-api/telegram-bot-api)
- PostgreSQL 17
- [goose](https://github.com/pressly/goose) — миграции
- Docker / docker compose

## Требования

- Go 1.26+
- Docker и docker compose
- goose:
  ```bash
  # ставится в $(go env GOPATH)/bin — эта папка должна быть в PATH
  go install github.com/pressly/goose/v3/cmd/goose@v3.26.0
  goose -version   # проверить, что установился
  ```

## Инициализация

1. Подтянуть зависимости:
   ```bash
   go mod download
   ```

2. Создать `.env` из шаблона:
   ```bash
   cp .env.example .env
   ```
   `.env` в git не коммитится — у каждого свой.

3. Заполнить значения в `.env`. Минимум — `BOT_TOKEN`, его выдаёт [@BotFather](https://t.me/BotFather).

## Запуск локально

1. Поднять Postgres:
   ```bash
   docker compose -f deployments/docker-compose.yml up -d db
   ```

2. Накатить миграции (goose сам читает `GOOSE_*` из `.env`):
   ```bash
   goose up
   ```

3. Запустить бота:
   ```bash
   go run ./cmd/bot
   ```

## Запуск через docker compose

В `.env` поменять хост базы на имя сервиса:

```
POSTGRES_HOST=db
GOOSE_DBSTRING=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
```

И поднять всё:

```bash
docker compose -f deployments/docker-compose.yml up -d --build
docker compose -f deployments/docker-compose.yml logs -f bot
```

Остановить:

```bash
docker compose -f deployments/docker-compose.yml down        # контейнеры
docker compose -f deployments/docker-compose.yml down -v     # + удалить данные базы
```

## Миграции

```bash
goose create <name> sql   # новая миграция в ./migrations
goose up                  # накатить все
goose down                # откатить последнюю
goose status              # статус
```
