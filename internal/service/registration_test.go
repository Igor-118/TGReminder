package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"tg-reminder/internal/domain"
	"tg-reminder/internal/service"
)

func TestRegistrationService_Start(t *testing.T) {
	tests := []struct {
		name string
		input *fakeUserRepo
		tgID int64
		want string
	}{
		{
			name: "пользователь уже зарегистрирован",
			input: &fakeUserRepo{
				user: &domain.User{Name: "Игорь"},
			},
			tgID: 50,
			want: "Игорь, ты уже зарегистрирован",
		},
		{
			name: "пользователя нет",
			input: &fakeUserRepo{
				err: domain.ErrNotFound,
			},
			tgID: 51,
			want: "Привет! Как тебя зовут?",
		},
		{
			name: "ошибка обёрнута",
			input: &fakeUserRepo{
				err: fmt.Errorf("get user: %w", domain.ErrNotFound),
			},
			tgID: 52,
			want: "Привет! Как тебя зовут?",
		},
		{
			name: "ошибка базы данных",
			input: &fakeUserRepo{
				err: errors.New("db down"),
			},
			tgID: 53,
			want: "Что-то пошло не так, попробуй позже",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			registration := service.NewRegistrationService(tt.input)

			got := registration.Start(ctx, tt.tgID)

			if tt.input.gotTgID != tt.tgID {
				t.Errorf("GetUserByTgID() получил ID %d, ожидали %d",
					tt.input.gotTgID, tt.tgID)
			}

			if got.Text != tt.want {
				t.Errorf("Start().Text = %q, ожидали %q", got.Text, tt.want)
			}
		})
	}
}
