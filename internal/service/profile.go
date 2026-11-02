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
	userRepo userGetter
}

func NewProfileService(userRepo userGetter) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

func (p *ProfileService) Show(ctx context.Context, tgID int64) Reply {
	user, err := p.userRepo.GetUserByTgID(ctx, tgID)
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
