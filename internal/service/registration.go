package service

import (
	"context"
	"fmt"
	"log/slog"
	"tg-reminder/internal/repository"
)

type RegistrationService struct {
	userRepo *repository.UserRepository
}

func NewRegistrationService(repository *repository.UserRepository) *RegistrationService{
	return &RegistrationService{userRepo: repository}
}

func (r *RegistrationService) HandleInput(ctx context.Context, userID int64, input string) Reply {
	user, err := r.userRepo.GetUserByTgID(ctx, userID)
	if err != nil {
		if err == repository.ErrNotFound {
			return Reply{Text: "Сначала зарегистрируйся — напиши /start"}
		}
		slog.Error("failed to get user",
			"userID", userID,
			"error", err)
		return Reply{Text: "Что-то пошло не так, попробуй позже"}
	}
	stroka := fmt.Sprintf("Привет, %v!", user.Name)
	return Reply{Text: stroka}
}
