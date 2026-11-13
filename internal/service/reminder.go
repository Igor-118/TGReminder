package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"tg-reminder/internal/domain"
	"tg-reminder/internal/repository"
	"time"
)

type ReminderService struct {
	userRepo     userGetter
	reminderRepo reminderGetter // *repository.ReminderRepository
}

type reminderGetter interface {
	ListUpcoming(ctx context.Context, userID int64, limit int) ([]domain.Reminder, error)
	Create(ctx context.Context, rem *domain.Reminder) error
}

func NewReminderService(userRepo userGetter, reminderRepo reminderGetter,
) *ReminderService {
	return &ReminderService{
		userRepo:     userRepo,
		reminderRepo: reminderRepo,
	}
}

const listLimit = 20

func (s *ReminderService) Add(ctx context.Context, tgID int64, args string) Reply {
	user, err := s.userRepo.GetUserByTgID(ctx, tgID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}
		slog.Error("failed to get user", "tgID", tgID, "error", err)
		return errorReply
	}
	loc, err := time.LoadLocation(user.TimeZone)
	if err != nil {
		slog.Error("не удалось получить часовой пояс", "userID", user.ID, "timeZone", user.TimeZone, "error", err)
		return errorReply
	}
	nowInZone := time.Now().In(loc)
	date, text, err := parseReminder(args, nowInZone, loc)
	if err != nil {
		if errors.Is(err, ErrTooLong) {
			return Reply{Text: "Текст слишком длинный"}
		}
		if errors.Is(err, ErrEmptyText) {
			return Reply{Text: "Пустой текст в напоминании"}
		}
		if errors.Is(err, ErrInPast) {
			return Reply{Text: "Нельзя создать напоминание в прошлом"}
		}
		if errors.Is(err, ErrBadFormat) {
			return Reply{
				Text: "Неверный формат. Используй: ДД.ММ.ГГГГ ЧЧ:ММ текст",
			}
		}

		slog.Error("не удалось считать строку", "userID", user.ID, "Text", text, "error", err)
		return errorReply
	}
	reminder := &domain.Reminder{
		UserID:   user.ID,
		Text:     text,
		RemindAt: date,
	}
	err = s.reminderRepo.Create(ctx, reminder)
	if err != nil {
		slog.Error("не удалось создать напоминание", "tgID", tgID, "error", err)
		return errorReply
	}
	return Reply{
		Text: fmt.Sprintf("Напоминание: %s", reminder.Text),
	}
}

func (s *ReminderService) List(ctx context.Context, tgID int64) Reply {
	user, err := s.userRepo.GetUserByTgID(ctx, tgID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}

		slog.Error("failed to get user", "tgID", tgID, "error", err)
		return errorReply
	}

	loc, err := time.LoadLocation(user.TimeZone)
	if err != nil {
		slog.Error("не удалось получить часовой пояс", "userID", user.ID, "timeZone", user.TimeZone, "error", err)
		return errorReply
	}

	reminders, err := s.reminderRepo.ListUpcoming(ctx, user.ID, listLimit)
	if err != nil {
		slog.Error("ошибка с листом напоминаний", "userID", user.ID, "error", err)
		return errorReply
	}
	if len(reminders) == 0 {
		return Reply{Text: "У тебя нет предстоящих напоминаний"}
	}

	var builder strings.Builder
	builder.WriteString("Ближайшие напоминания:\n")
	for i, reminder := range reminders {
		remindAt := reminder.RemindAt.In(loc)
		if reminder.Text == "" {
			fmt.Fprintf(&builder, "%d. Без описания — %s\n", i+1,
				remindAt.Format("02.01.2006 15:04"),
			)
		} else {
			fmt.Fprintf(&builder, "%d. %s — %s\n", i+1,
				reminder.Text,
				remindAt.Format("02.01.2006 15:04"),
			)
		}

	}

	return Reply{Text: builder.String()}
}
