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

func TestData(t *testing.T) {
	tests := []struct {
		name  string
		input *fakeUserRepo
		want  []string
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
				gotTgID: 50,
			},
			want: []string{"Бурмалда", "228", "Europe/Moscow", "Сейчас у тебя:"},
		},
		{
			name: "пользователя нет",
			input: &fakeUserRepo{
				user:    &domain.User{},
				err:     repository.ErrNotFound,
				gotTgID: 51,
			},
			want: []string{
				"Сначала зарегистрируйся — напиши /start",
			},
		},
		{
			name: "ошибка обёрнута в лаваш",
			input: &fakeUserRepo{
				user:    &domain.User{},
				err:     fmt.Errorf("get user: %w", repository.ErrNotFound),
				gotTgID: 52,
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
				gotTgID: 53,
			},
			want: []string{
				"Что-то пошло не так, попробуй позже",
			},
		},
		{
			name: "неверный пояс чаовой",
			input: &fakeUserRepo{
				user: &domain.User{
					ID:       101,
					Name:     "kuku",
					Age:      228,
					TimeZone: "Mars/Olympus",
				},
				err:     nil,
				gotTgID: 54,
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
			got := profile.Show(ctx, 50)
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
