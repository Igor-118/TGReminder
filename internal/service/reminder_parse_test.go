package service

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestReminderParse(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(
		2027, 1, 1,
		12, 0, 0, 0,
		loc,
	)
	tests := []struct {
		name     string
		args     string
		wantErr  error
		wantText string
	}{
		{
			name:     "корректное напоминание",
			args:     "02.01.2027 15:00 вайбкодинг отстой!",
			wantText: "вайбкодинг отстой!",
		},
		{
			name:    "пустая строка",
			args:    "",
			wantErr: ErrEmptyText,
		},
		{
			name:    "нет текста",
			args:    "02.01.2027 15:00",
			wantErr: ErrEmptyText,
		},
		{
			name:    "неправильная дата",
			args:    "32.01.2027 15:00 Купить молоко",
			wantErr: ErrBadFormat,
		},
		{
			name:    "неправильный месяц",
			args:    "32.13.2027 15:00 Купить молоко",
			wantErr: ErrBadFormat,
		},
		{
			name:    "неправильное время",
			args:    "32.01.2027 25:00 Купить молоко",
			wantErr: ErrBadFormat,
		},
		{
			name:    "неправильные минуты",
			args:    "02.01.2027 15:99 Купить молоко",
			wantErr: ErrBadFormat,
		},
		{
			name:    "дата в прошлом",
			args:    "01.01.2027 10:00 Купить молоко",
			wantErr: ErrInPast,
		},
		{
			name:    "слишком длинный текст",
			args:    strings.Repeat("a", maxSymbols+1),
			wantErr: ErrTooLong,
		},
		{
			name:    "длинна сообщения меньше 2 слов",
			args:    "02.01.2028",
			wantErr: ErrBadFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDate, gotText, gotErr := parseReminder(
				tt.args,
				now,
				loc,
			)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatalf(
					"parseReminder() error = %v, want %v",
					gotErr,
					tt.wantErr,
				)
			}
			if tt.wantErr == nil {
				if gotText != tt.wantText {
					t.Errorf(
						"parseReminder() text = %q, want %q",
						gotText,
						tt.wantText,
					)
				}
				wantDate := time.Date(
					2027, 1, 2,
					15, 0, 0, 0,
					loc,
				)
				if !gotDate.Equal(wantDate) {
					t.Errorf(
						"parseReminder() date = %v, want %v",
						gotDate,
						wantDate,
					)
				}
			}
		})
	}
}
