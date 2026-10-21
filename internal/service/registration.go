package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tg-reminder/internal/repository"
	"tg-reminder/internal/domain"
)

type RegistrationService struct {
	userInterface UserInterface
}

type UserInterface interface{
	 GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error)
}

func NewRegistrationService(repository *repository.UserRepository) *RegistrationService {
	return &RegistrationService{userInterface: repository}
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
