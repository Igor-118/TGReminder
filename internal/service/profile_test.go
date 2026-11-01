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

type fakeUserRepo struct {
	user    *domain.User // какого пользователя вернуть
	err     error        // какую ошибку вернуть
	gotTgID int64        // сюда метод запишет, с каким tgID его вызвали
}

func (f *fakeUserRepo) GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	f.gotTgID = tgID
	return f.user, f.err
}

func TestProfileService_Show(t *testing.T) {
	tests := []struct {
		name  string
		input *fakeUserRepo
		want  []string
		tgID  int64
	}{
		{
			name: "Пользователь существует",
			input: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
				err:     nil,
			},
			want: []string{"Бурмалда", "228", "Europe/Moscow", "Сейчас у тебя:"},
		},
		{
			name: "пользователя нет",
			input: &fakeUserRepo{
				user:    &domain.User{},
				err:     repository.ErrNotFound,
			},
			want: []string{
				"Сначала зарегистрируйся — напиши /start",
			},
		},
		{
			name: "ошибка обёрнута",
			input: &fakeUserRepo{
				user:    &domain.User{},
				err:     fmt.Errorf("get user: %w", repository.ErrNotFound),
			},
			want: []string{
				"Сначала зарегистрируйся — напиши /start",
			},
		},
		{
			name: "ошибка Базы данных",
			input: &fakeUserRepo{
				user:    &domain.User{},
				err:     errors.New("db down, База данных вниз"),
			},
			want: []string{
				"Что-то пошло не так, попробуй позже",
			},
		},
		{
			name: "неверный пояс чаовой",
			input: &fakeUserRepo{
				user: &domain.User{
					TimeZone: "Mars/Olympus",
				},
				err:     nil,
			},
			want: []string{
				"Что-то пошло не так, попробуй позже",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			profile := service.NewProfileService(tt.input)
			got := profile.Show(ctx, tt.tgID)
			if tt.input.gotTgID != tt.tgID {
				t.Errorf(
					"GetUserByTgID() получил ID %d, ожидали %d",
					tt.input.gotTgID,
					tt.tgID,
				)
			}
			for _, expected := range tt.want {
				if !strings.Contains(got.Text, expected) {
					t.Errorf(
						"Show().Text = %q, ожидали наличие %q",
						got.Text, expected,
					)
				}
			}

		})
	}
}
