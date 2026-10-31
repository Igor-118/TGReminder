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
	"time"
)

type fakeReminderUserRepo struct {
	user    *domain.User // какого пользователя вернуть
	err     error        // какую ошибку вернуть
	gotTgID int64        // сюда метод запишет, с каким tgID его вызвали
}

type fakeReminderRepo struct {
	reminders []domain.Reminder // список событий у возвращенного пользователя
	err       error             // какую ошибку вернуть
	gotUserID int64             // сюда метод запишет, с каким tgID его вызвали
	gotLimit  int
	called    bool
}

func (f *fakeReminderUserRepo) GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	f.gotTgID = tgID
	return f.user, f.err
}

func (f *fakeReminderRepo) ListUpcoming(ctx context.Context, userID int64, limit int) ([]domain.Reminder, error) {
	f.called = true
	f.gotUserID = userID
	f.gotLimit = limit
	return f.reminders, f.err
}

func TestDataReminder(t *testing.T) {
	tests := []struct {
		name               string
		inputUser          *fakeReminderUserRepo
		inputReminder      *fakeReminderRepo
		tgID               int64
		wantUserID         int64
		wantReminderCalled bool
		want               []string
	}{
		{
			name: "список дел",
			inputUser: &fakeReminderUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
				err:     nil,
				gotTgID: 50,
			},
			inputReminder: &fakeReminderRepo{
				reminders: []domain.Reminder{
					{
						Text:     "Сходить в магазин",
						RemindAt: time.Date(2026, 11, 1, 15, 0, 0, 0, time.UTC),
					},
					{
						Text:     "вернуться с магазина",
						RemindAt: time.Date(2026, 11, 2, 16, 0, 0, 0, time.UTC),
					},
				},
			},
			tgID:               50,
			wantUserID:         100,
			wantReminderCalled: true,
			want:               []string{"Сходить в магазин", "вернуться с магазина"},
		},
		{
			name: "пользователя нет",
			inputUser: &fakeReminderUserRepo{
				user: &domain.User{},
				err:  repository.ErrNotFound,
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          51,
			want:          []string{"Сначала зарегистрируйся — напиши /start"},
		},
		{
			name: "ошибка обёрнута",
			inputUser: &fakeReminderUserRepo{
				err: fmt.Errorf("get user: %w", repository.ErrNotFound),
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          52,
			want:          []string{"Сначала зарегистрируйся — напиши /start"},
		},
		{
			name: "ошибка Базы данных",
			inputUser: &fakeReminderUserRepo{
				err: errors.New("db down, База данных вниз"),
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          53,
			want:          []string{"Что-то пошло не так, попробуй позже"},
		},
		{
			name: "неверный пояс чаовой",
			inputUser: &fakeReminderUserRepo{
				user: &domain.User{
					ID:       100,
					TimeZone: "Mars/Olympus",
				},
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          54,
			want:          []string{"Что-то пошло не так, попробуй позже"},
		},
		{
			name: "ошибка с листом напоминаний",
			inputUser: &fakeReminderUserRepo{
				user: &domain.User{
					ID:       100,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder: &fakeReminderRepo{
				err: errors.New("db down"),
			},
			tgID:               55,
			wantUserID:         100,
			wantReminderCalled: true,
			want:               []string{"Что-то пошло не так, попробуй позже"},
		},
		{
			name: "пустой список напоминаний",
			inputUser: &fakeReminderUserRepo{
				user: &domain.User{
					ID:       100,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder: &fakeReminderRepo{
				reminders: []domain.Reminder{},
			},
			tgID:               56,
			wantUserID:         100,
			wantReminderCalled: true,
			want: []string{
				"У тебя нет предстоящих напоминаний",
			},
		},
		{
			name: "напоминание без описания",
			inputUser: &fakeReminderUserRepo{
				user: &domain.User{
					ID:       100,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder: &fakeReminderRepo{
				reminders: []domain.Reminder{
					{
						Text:     "",
						RemindAt: time.Date(2026, 11, 3, 15, 0, 0, 0, time.UTC),
					},
				},
			},
			tgID:               57,
			wantUserID:         100,
			wantReminderCalled: true,
			want: []string{
				"Без описания",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			reminderService := service.NewReminderService(tt.inputUser, tt.inputReminder)
			got := reminderService.List(ctx, tt.tgID)
			if tt.inputUser.gotTgID != tt.tgID {
				t.Errorf(
					"GetUserByTgID() получил ID %d, ожидали %d",
					tt.inputUser.gotTgID,
					tt.tgID,
				)
			}

			// Проверяем, вызывался ли репозиторий напоминаний.
			if tt.inputReminder.called != tt.wantReminderCalled {
				t.Errorf(
					"ListUpcoming() вызван: %v, ожидали: %v",
					tt.inputReminder.called,
					tt.wantReminderCalled,
				)
			}

			// Если репозиторий напоминаний должен был вызваться,
			// проверяем переданный внутренний ID пользователя.
			if tt.wantReminderCalled {
				if tt.inputReminder.gotUserID != tt.wantUserID {
					t.Errorf(
						"ListUpcoming() получил ID %d, ожидали %d",
						tt.inputReminder.gotUserID,
						tt.wantUserID,
					)
				}

				if tt.inputReminder.gotLimit != 20 {
					t.Errorf(
						"ListUpcoming() получил limit %d, ожидали 20",
						tt.inputReminder.gotLimit,
					)
				}
			}

			// Проверяем текст ответа.
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
