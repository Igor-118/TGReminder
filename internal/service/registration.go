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
	userInterface userGetter
}

func NewRegistrationService(userRepo userGetter) *RegistrationService {
	return &RegistrationService{userInterface: userRepo}
}

type userGetter interface {
	GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error)
}

func (r *RegistrationService) HandleInput(ctx context.Context, userID int64, input string) Reply {
	slog.Info("Starting HandleInput", "userID", userID)
	user, err := r.userInterface.GetUserByTgID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}
		slog.Error("failed to get user",
			"userID", userID,
			"error", err)
		return errorReply
	}
	slog.Info("лог что поиск успешно завершён", "userID", userID)
	return Reply{Text: fmt.Sprintf("Привет, %v!", user.Name)}
}
