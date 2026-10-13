package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"tg-reminder/internal/repository"
	"time"
)

type ReminderService struct {
	userRepo     *repository.UserRepository
	reminderRepo *repository.ReminderRepository
}

func NewReminderService(userRepo *repository.UserRepository, reminderRepo *repository.ReminderRepository,
) *ReminderService {
	return &ReminderService{
		userRepo:     userRepo,
		reminderRepo: reminderRepo,
	}
}

const listLimit = 20

func (s *ReminderService) List(ctx context.Context, userID int64) Reply {
	user, err := s.userRepo.GetUserByTgID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}

		slog.Error("failed to get user", "userID", userID, "error", err)
		return errorReply
	}

	loc, err := time.LoadLocation(user.TimeZone)
	if err != nil {
		slog.Error("не удалось получить часовой пояс", "userID", userID, "timeZone", user.TimeZone, "error", err)
		return errorReply
	}

	reminders, err := s.reminderRepo.ListUpcoming(ctx, user.ID, listLimit)
	if err != nil {
		slog.Error("ошибка с листом напоминаний", "userID", userID, "error", err)
		return errorReply
	}
	if len(reminders) == 0 {
		return Reply{Text: "У тебя нет предстоящих напоминаний"}
	}

	var builder strings.Builder
	builder.WriteString("Ближайшие напоминания:\n")
	for i, reminder := range reminders {
		remindeAt := reminder.RemindAt.In(loc)
		if reminder.Text == "" {
			fmt.Fprintf(&builder, "%d. Без описания — %s\n", i+1,
				remindeAt.Format("02.01.2006 15:04"),
			)
		} else {
			fmt.Fprintf(&builder, "%d. %s — %s\n", i+1,
				reminder.Text,
				remindeAt.Format("02.01.2006 15:04"),
			)
		}

	}

	return Reply{Text: builder.String()}
}
