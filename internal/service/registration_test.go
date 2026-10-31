package service_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
		want  []string
	}{
		{
			name: "приветствие пользователя",
			input: &fakeUserInterface{
				user: &domain.User{
					Name: "Игорь",
				},
				err:     nil,
				gotTgID: 10,
			},
			want: []string{"Привет, Игорь!"},
		},
		{
			name: "пользователя нет",
			input: &fakeUserInterface{
				user: &domain.User{},
				err:  repository.ErrNotFound,
			},
			want: []string{"Сначала зарегистрируйся — напиши /start"},
		},
		{
			name: "ошибка обёрнута",
			input: &fakeUserInterface{
				user:    &domain.User{},
				err:     fmt.Errorf("get user: %w", repository.ErrNotFound),
				gotTgID: 11,
			},
			want: []string{"Сначала зарегистрируйся — напиши /start"},
		},
		{
			name: "ошибка Базы Данных",
			input: &fakeUserInterface{
				user:    &domain.User{},
				err:     errors.New("db down, База данных вниз"),
				gotTgID: 12,
			},
			want: []string{"Что-то пошло не так, попробуй позже"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			registration := service.NewRegistrationService(tt.input)
			got := registration.HandleInput(ctx, 10, "kek")
			for _, expected := range tt.want {
				if !strings.Contains(got.Text, expected) {
					t.Errorf(
						"Show().Text = %q, ожидали. А получили: %q",
						got.Text, expected,
					)
				}
			}
		})
	}
}
