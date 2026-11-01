package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tg-reminder/internal/domain"
	"tg-reminder/internal/repository"
)

type RegistrationService struct {
	userRepo userGetter
}

func NewRegistrationService(userRepo userGetter) *RegistrationService {
	return &RegistrationService{userRepo: userRepo}
}

type userGetter interface {
	GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error)
}

func (r *RegistrationService) HandleInput(ctx context.Context, tgID int64, input string) Reply {
	slog.Info("Starting HandleInput", "userID", tgID)
	user, err := r.userRepo.GetUserByTgID(ctx, tgID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}
		slog.Error("failed to get user",
			"userID", tgID,
			"error", err)
		return errorReply
	}
	slog.Info("лог что поиск успешно завершён", "userID", tgID)
	return Reply{Text: fmt.Sprintf("Привет, %v!", user.Name)}
}
