package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tg-reminder/internal/repository"
	"time"
)

type ProfileService struct {
	userRepo *repository.UserRepository
}

func NewProfileService(userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

func (p *ProfileService) Show(ctx context.Context, userID int64) Reply {
	user, err := p.userRepo.GetUserByTgID(ctx, userID)
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

	nowInZone := time.Now().In(loc)
	return Reply{
		Text: fmt.Sprintf(
			"Твой профиль:\nИмя: %s\nВозраст: %d\nЧасовой пояс: %s\nСейчас у тебя: %s",
			user.Name,
			user.Age,
			user.TimeZone,
			nowInZone.Format("15:04"),
		),
	}
}
