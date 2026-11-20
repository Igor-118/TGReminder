package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"tg-reminder/internal/domain"
)

var timezoneOptions = []Option{
	{Label: "Калининград", Value: "Europe/Kaliningrad"},
	{Label: "Москва", Value: "Europe/Moscow"},
	{Label: "Самара", Value: "Europe/Samara"},
	{Label: "Екатеринбург", Value: "Asia/Yekaterinburg"},
	{Label: "Омск", Value: "Asia/Omsk"},
	{Label: "Новосибирск", Value: "Asia/Novosibirsk"},
	{Label: "Красноярск", Value: "Asia/Krasnoyarsk"},
	{Label: "Иркутск", Value: "Asia/Irkutsk"},
	{Label: "Якутск", Value: "Asia/Yakutsk"},
	{Label: "Владивосток", Value: "Asia/Vladivostok"},
	{Label: "Магадан", Value: "Asia/Magadan"},
	{Label: "Камчатка", Value: "Asia/Kamchatka"},
}

type RegistrationService struct {
	userRepo userStore
	states   *StateStore
}

func NewRegistrationService(userRepo userStore) *RegistrationService {
	return &RegistrationService{
		userRepo: userRepo,
		states:   NewStateStore(),
	}
}

type userGetter interface {
	GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error)
}

type userStore interface {
	GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
}

func (s *RegistrationService) Start(ctx context.Context, tgID int64) Reply {
	slog.Info("Starting Registration", "tgID", tgID)

	user, err := s.userRepo.GetUserByTgID(ctx, tgID)
	if err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			slog.Error("failed to get user",
				"tgID", tgID,
				"error", err)
			return errorReply
		}

		s.states.Set(tgID, &registration{
			Step: stepName,
		})

		return Reply{Text: "Привет! Как тебя зовут?"}
	}

	return Reply{Text: fmt.Sprintf("%s, ты уже зарегистрирован", user.Name)}
}

func (r *RegistrationService) HandleInput(ctx context.Context, tgID int64, input string) Reply {
	slog.Info("Starting HandleInput", "tgID", tgID)

	reg, ok := r.states.Get(tgID)
	if !ok {
		return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
	}

	switch reg.Step {
	case stepName:
		name := strings.TrimSpace(input)

		if name == "" {
			return Reply{Text: "Имя не может быть пустым"}
		}

		if utf8.RuneCountInString(name) > 64 {
			return Reply{Text: "Имя не должно быть длиннее 64 символов"}
		}

		reg.Name = name
		reg.Step = stepAge
		r.states.Set(tgID, reg)

		return Reply{Text: "Сколько тебе лет?"}

	case stepAge:
		age, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || age < 1 || age > 120 {
			return Reply{Text: "Возраст должен быть целым числом от 1 до 120"}
		}

		reg.Age = age
		reg.Step = stepTimezone
		r.states.Set(tgID, reg)

		return Reply{
			Text:    "Выбери часовой пояс:",
			Options: timezoneOptions,
		}

	case stepTimezone:
		for _, option := range timezoneOptions {
			if option.Value == input {
				_, err := time.LoadLocation(input)
				if err != nil {
					slog.Error("failed to load timezone",
						"timezone", input,
						"error", err)
					return errorReply
				}

				user := &domain.User{
					TgID:     tgID,
					Name:     reg.Name,
					Age:      reg.Age,
					TimeZone: input,
				}

				err = r.userRepo.Create(ctx, user)
				if err != nil {
					slog.Error("failed to create user",
						"tgID", tgID,
						"error", err)
					return errorReply
				}

				r.states.Delete(tgID)

				return Reply{Text: "Регистрация завершена!"}
			}
		}

		return Reply{Text: "Выбери часовой пояс из списка"}

	default:
		return errorReply
	}
}
