# 2. Проверка, зарегистрирован ли пользователь

**Зависит от:** [0](00-tg-id.md), [1](01-config-db.md)

## Зачем

Без часового пояса бот ничего не может: не поставишь напоминание, пока не знаешь, когда у человека 10 утра. Поэтому каждое сообщение сначала проверяется: зарегистрирован ли пользователь. Если нет, бот просит написать `/start`, саму регистрацию делаем в следующих подзадачах.

## Что сделать

### `internal/domain/user.go`

```go
type User struct {
	ID       int64
	TgID     int64
	Name     string
	Age      int
	TimeZone string
}
```

### `internal/repository/repository.go`

```go
var ErrNotFound = errors.New("not found")

type UserRepository interface {
	// Возвращает ErrNotFound, если пользователя нет.
	GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error)
}
```

Возвращаем пользователя целиком, а не просто `bool`: тогда в сообщениях можно обращаться к нему по имени.

### `internal/repository/postgres/user.go`

Реализация `UserRepository` на pgx. Если строки нет, pgx вернёт `pgx.ErrNoRows`. Её нужно превратить в `repository.ErrNotFound`, чтобы сервис не зависел от pgx:

```go
if errors.Is(err, pgx.ErrNoRows) {
	return nil, repository.ErrNotFound
}
```

### `internal/service/reply.go`

Ответ сервиса. Сервис не знает про Telegram, поэтому возвращает вот такую структуру, а бот сам решает, как её показать.

```go
// Reply — что показать пользователю.
type Reply struct {
	Text    string
	Options []Option // варианты на выбор; пусто — просто текст
}

// Option — один вариант ответа.
type Option struct {
	Label string // что видит пользователь: "Москва"
	Value string // что вернётся в сервис при выборе: "Europe/Moscow"
}
```

`Options` в этой подзадаче не используются, они понадобятся в [3](03-registration-service.md) для выбора часового пояса.

### `internal/service/registration.go`

```go
type RegistrationService struct {
	users repository.UserRepository
}

// Каждое сообщение пользователя сначала идёт сюда.
// user != nil — пользователь зарегистрирован, бот обрабатывает сообщение дальше сам.
// user == nil — нужно отправить пользователю reply и больше ничего не делать.
func (s *RegistrationService) HandleInput(ctx context.Context, userID int64, input string) (Reply, *domain.User, error) {
	// GetUserByTgID:
	//   нашёлся          → вернуть пользователя
	//   ErrNotFound      → Reply{Text: "Сначала зарегистрируйся — напиши /start"}
	//   другая ошибка    → вернуть ошибку
}
```

### `internal/bot/handlers.go`

На каждое сообщение:
1. вызвать `HandleInput`;
2. если `user == nil`, отправить `reply.Text` и выйти;
3. иначе обработать как обычно, пока это эхо: `"<Имя>, ты написал: <текст>"`.

Если `HandleInput` вернул ошибку, залогировать её и ответить пользователю «Что-то пошло не так, попробуй позже». Бот при этом падать не должен.

### `internal/app/app.go`

Склеить: пул → `postgres`-репозиторий → `RegistrationService` → бот.

## Готово, когда

- [ ] незарегистрированный пользователь на любое сообщение получает «напиши /start»
- [ ] если руками вставить себя в `users`, бот отвечает эхом и обращается по имени
- [ ] если остановить Postgres, бот не падает, а отвечает «что-то пошло не так»
- [ ] `internal/service` не импортирует `tgbotapi` и pgx:
  ```bash
  go list -deps ./internal/service | grep -E 'telegram|pgx'   # должно быть пусто
  ```

## Не входит

- `/start` и сама регистрация — [3](03-registration-service.md), [4](04-registration-bot.md)
