package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"tg-reminder/internal/repository"
)

type RegistrationService struct {
	userRepo *repository.UserRepository
}

func NewRegistrationService(repository *repository.UserRepository) *RegistrationService {
	return &RegistrationService{userRepo: repository}
}

func (r *RegistrationService) HandleInput(ctx context.Context, userID int64, input string) Reply {
	slog.Info("Starting HandleInput", "userID", userID)
	user, err := r.userRepo.GetUserByTgID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}
		slog.Error("failed to get user",
			"userID", userID,
			"error", err)
		return Reply{Text: "Что-то пошло не так, попробуй позже"}
	}
	slog.Info("лог что поиск успешно завершён", userID)
	return Reply{Text: fmt.Sprintf("Привет, %v!", user.Name)}
}
