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

type fakeUserInterface struct {
	user    *domain.User // какого пользователя вернуть
	err     error        // какую ошибку вернуть
	gotTgID int64        // сюда метод запишет, с каким tgID его вызвали
}

func (f *fakeUserInterface) GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	f.gotTgID = tgID
	return f.user, f.err
}

func TestDataInterface(t *testing.T) {
	tests := []struct {
		name  string
		input *fakeUserInterface
		want  string
		tgID  int64
	}{
		{
			name: "приветствие пользователя",
			input: &fakeUserInterface{
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
			input: &fakeUserInterface{
				user: &domain.User{},
				err:  repository.ErrNotFound,
			},
			tgID: 51,
			want: "Сначала зарегистрируйся — напиши /start",
		},
		{
			name: "ошибка обёрнута",
			input: &fakeUserInterface{
				user: &domain.User{},
				err:  fmt.Errorf("get user: %w", repository.ErrNotFound),
			},
			tgID: 52,
			want: "Сначала зарегистрируйся — напиши /start",
		},
		{
			name: "ошибка Базы Данных",
			input: &fakeUserInterface{
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
					"HandleInput().Text = %q, ожидали. А получили: %q",
					tt.want, got.Text,
				)
			}

		})
	}
}
