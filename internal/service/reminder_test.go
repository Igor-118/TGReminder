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

func TestReminderService_Add(t *testing.T) {
	tests := []struct {
		name               string
		inputUser          *fakeUserRepo
		inputReminder      *fakeReminderRepo
		tgID               int64
		wantUserID         int64
		wantReminderCalled bool
		args               string
		want               string
	}{
		{
			name: "успешно добавлено",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               50,
			wantUserID:         100,
			wantReminderCalled: true,
			args:               "12.12.2027 15:00 Сходить в магазин",
			want:               "Напоминание: Сходить в магазин",
		},
		{
			name: "пустой текст в напоминании",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               51,
			wantReminderCalled: false,
			args:               "12.12.2027 15:00",
			want:               "Пустой текст в напоминании",
		},
		{
			name: "пользователя нет",
			inputUser: &fakeUserRepo{
				user: nil,
				err:  repository.ErrNotFound,
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          52,
			want:          "Сначала зарегистрируйся — напиши /start",
		},
		{
			name: "ошибка обёрнута",
			inputUser: &fakeUserRepo{
				err: fmt.Errorf("get user: %w", repository.ErrNotFound),
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          52,
			want:          "Сначала зарегистрируйся — напиши /start",
		},
		{
			name: "ошибка Базы данных",
			inputUser: &fakeUserRepo{
				err: errors.New("db down, База данных вниз"),
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          53,
			want:          "Что-то пошло не так, попробуй позже",
		},
		{
			name: "неверный пояс чаовой",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					TimeZone: "Mars/Olympus",
				},
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          54,
			want:          "Что-то пошло не так, попробуй позже",
		},
		{
			name: "неверный формат числа",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               51,
			wantReminderCalled: false,
			args:               "32.12.2027 15:00 текст",
			want:               "Неверный формат. Используй: ДД.ММ.ГГГГ ЧЧ:ММ текст",
		},
		{
			name: "неверный формат месяца",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               51,
			wantReminderCalled: false,
			args:               "12.121.2027 15:00 текст",
			want:               "Неверный формат. Используй: ДД.ММ.ГГГГ ЧЧ:ММ текст",
		},
		{
			name: "неверный формат года",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               51,
			wantReminderCalled: false,
			args:               "12.12.20271 15:00 текст",
			want:               "Неверный формат. Используй: ДД.ММ.ГГГГ ЧЧ:ММ текст",
		},
		{
			name: "неверный формат времени",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               51,
			wantReminderCalled: false,
			args:               "12.12.2027 25:00 текст",
			want:               "Неверный формат. Используй: ДД.ММ.ГГГГ ЧЧ:ММ текст",
		},
		{
			name: "дата в прошлом",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
			},
			inputReminder:      &fakeReminderRepo{},
			tgID:               55,
			args:               "01.01.2026 15:00 Сходить в магазин",
			wantReminderCalled: false,
			want:               "Нельзя создать напоминание в прошлом",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			reminderService := service.NewReminderService(tt.inputUser, tt.inputReminder)
			got := reminderService.Add(ctx, tt.tgID, tt.args)
			if tt.inputUser.gotTgID != tt.tgID {
				t.Errorf(
					"GetUserByTgID() получил ID %d, ожидали %d",
					tt.inputUser.gotTgID,
					tt.tgID,
				)
			}
			if tt.inputReminder.called != tt.wantReminderCalled {
				t.Errorf(
					"Create() called = %v, ожидали %v",
					tt.inputReminder.called,
					tt.wantReminderCalled,
				)
			}
			if tt.wantReminderCalled && tt.inputReminder.gotUserID != tt.wantUserID {
				t.Errorf(
					"Create() получил userID %d, ожидали %d",
					tt.inputReminder.gotUserID,
					tt.wantUserID,
				)
			}
			if got.Text != tt.want {
				t.Errorf(
					"Add().Text = %q, ожидали наличие %q",
					got.Text, tt.want,
				)
			}

		})
	}
}

func TestReminderService_List(t *testing.T) {
	tests := []struct {
		name               string
		inputUser          *fakeUserRepo
		inputReminder      *fakeReminderRepo
		tgID               int64
		wantUserID         int64
		wantReminderCalled bool
		want               []string
	}{
		{
			name: "успешное получение списка дел",
			inputUser: &fakeUserRepo{
				user: &domain.User{
					ID:       100,
					Name:     "Бурмалда",
					Age:      228,
					TimeZone: "Europe/Moscow",
				},
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
			inputUser: &fakeUserRepo{
				user: nil,
				err:  repository.ErrNotFound,
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          51,
			want:          []string{"Сначала зарегистрируйся — напиши /start"},
		},
		{
			name: "ошибка обёрнута",
			inputUser: &fakeUserRepo{
				err: fmt.Errorf("get user: %w", repository.ErrNotFound),
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          52,
			want:          []string{"Сначала зарегистрируйся — напиши /start"},
		},
		{
			name: "ошибка Базы данных",
			inputUser: &fakeUserRepo{
				err: errors.New("db down, База данных вниз"),
			},
			inputReminder: &fakeReminderRepo{},
			tgID:          53,
			want:          []string{"Что-то пошло не так, попробуй позже"},
		},
		{
			name: "неверный пояс чаовой",
			inputUser: &fakeUserRepo{
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
			inputUser: &fakeUserRepo{
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
			inputUser: &fakeUserRepo{
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
			inputUser: &fakeUserRepo{
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
						"Add().Text = %q, ожидали наличие %q",
						got.Text, expected,
					)
				}
			}
		})
	}
}
