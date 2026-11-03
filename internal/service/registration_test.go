package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"tg-reminder/internal/domain"
	"tg-reminder/internal/repository"
	"tg-reminder/internal/service"
)

func TestRegistrationService_HandleInput(t *testing.T) {
	tests := []struct {
		name  string
		input *fakeUserRepo
		want  string
		tgID  int64
	}{
		{
			name: "приветствие пользователя",
			input: &fakeUserRepo{
				user: &domain.User{
					Name: "Игорь",
				},
				err: nil,
			},
			tgID: 50,
			want: "Привет, Игорь!",
		},
		{
			name: "пользователя нет",
			input: &fakeUserRepo{
				user: &domain.User{},
				err:  repository.ErrNotFound,
			},
			tgID: 51,
			want: "Сначала зарегистрируйся — напиши /start",
		},
		{
			name: "ошибка обёрнута",
			input: &fakeUserRepo{
				user: &domain.User{},
				err:  fmt.Errorf("get user: %w", repository.ErrNotFound),
			},
			tgID: 52,
			want: "Сначала зарегистрируйся — напиши /start",
		},
		{
			name: "ошибка Базы Данных",
			input: &fakeUserRepo{
				user: &domain.User{},
				err:  errors.New("db down, База данных вниз"),
			},
			tgID: 53,
			want: "Что-то пошло не так, попробуй позже",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			registration := service.NewRegistrationService(tt.input)
			got := registration.HandleInput(ctx, tt.tgID, "kek")
			if tt.input.gotTgID != tt.tgID {
				t.Errorf(
					"GetUserByTgID() получил ID %d, ожидали %d",
					tt.input.gotTgID,
					tt.tgID,
				)
			}
			if got.Text != tt.want {
				t.Errorf(
					"HandleInput().Text = %q, ожидали %q", got.Text, tt.want,
				)
			}

		})
	}
}
