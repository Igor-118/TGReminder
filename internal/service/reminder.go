package service

import (
	"context"
	"errors"
	"log/slog"
	"tg-reminder/internal/repository"
	"fmt"
	"strings"
	"time"
)

type ReminderService struct {
	userRepo *repository.UserRepository
	reminderRepo    *repository.ReminderRepository
}

func NewReminderService(userRepo *repository.UserRepository, reminderRepo *repository.ReminderRepository,
) *ReminderService {
	return &ReminderService{
		userRepo: userRepo,
		reminderRepo:    reminderRepo,
	}
}

var errorsReply = Reply{Text: "Что-то пошло не так, попробуй позже"}
const listLimit = 20

func (s *ReminderService) List(ctx context.Context, userID int64) Reply {
	user, err := s.userRepo.GetUserByTgID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}

		slog.Error("failed to get user", "userID", userID, "error", err)
		return errorsReply
	}

	loc, err := time.LoadLocation(user.TimeZone)
	if err != nil {
		slog.Error("не удалось получить часовой пояс", "userID", userID, "timeZone", user.TimeZone, "error", err)
		return errorsReply
	}

	reminders, err := s.reminderRepo.ListUpcoming(ctx, user.ID, listLimit)
	if err != nil {
		slog.Error("ошибка с листом напоминаний", "userID", userID, "error", err)
		return Reply{Text: "не получилось получить лист напоминаний"}
	}
	if len(reminders) == 0 {
		return Reply{Text: "У тебя нет предстоящих напоминаний"}
	}

	var builder strings.Builder
	builder.WriteString("Ближайшие напоминания:\n")
	for i, reminder := range reminders {
		remindeAt := reminder.RemindAt.In(loc)
		fmt.Fprintf(&builder, "%d. %s — %s\n", i+1,
			reminder.Text,
			remindeAt.Format("02.01.2006 15:04"),
		)
	}

	return Reply{Text: builder.String()}
}
