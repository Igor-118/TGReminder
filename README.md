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
- make
- [golangci-lint](https://golangci-lint.run/) v2 — только для `make lint`

## Инициализация

1. Подтянуть зависимости:
   ```bash
   go mod download
   ```

2. Положить токен бота в `.secrets` (в git не коммитится). Токен выдаёт [@BotFather](https://t.me/BotFather):
   ```bash
   echo 'BOT_TOKEN=123456:ABC...' > .secrets
   ```

`.env` руками не редактируется — его собирает `make` из `env.example` и `.secrets`.
Чтобы поменять настройки, правь `env.example`, а секреты — `.secrets`: значения из `.secrets` перекрывают одноимённые из `env.example`.

## Запуск локально

Бот запускается на хосте, Postgres — в докере:

```bash
make run-local
```

Что делает команда:

1. `make local-env` — собирает `.env` с хостом базы `127.0.0.1`;
2. `make db-up` — поднимает Postgres в докере и ждёт, пока он начнёт принимать подключения;
3. `make migrate-up` — накатывает миграции;
4. `go run ./cmd/bot` — запускает бота.

Остановить базу (данные сохраняются):

```bash
make db-down
```

## Запуск через docker compose

Бот и Postgres — оба в докере:

```bash
make run
docker compose --env-file .env -f deployments/docker-compose.yml logs -f bot   # логи бота
```

`make run` собирает `.env` для докера, пересобирает образ и поднимает контейнеры. Миграции накатываются при старте контейнера бота (`entrypoint.sh`).

Остановить:

```bash
make stop                                                       # контейнеры
docker compose --env-file .env -f deployments/docker-compose.yml down -v        # + удалить данные базы
```

## Миграции

Перед командами нужен `.env` для локального запуска (`make local-env`) и поднятая база (`make db-up`).

```bash
make migrate-up       # накатить все
make migrate-down     # откатить последнюю
make migrate-status   # статус
```

Новая миграция:

```bash
go run github.com/pressly/goose/v3/cmd/goose@v3.26.0 -dir migrations create <name> sql
```

## Тесты и линтер

```bash
make test   # тесты + отчёт о покрытии в coverage.html
make lint   # golangci-lint
```
