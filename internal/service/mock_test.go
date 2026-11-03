package service_test

import (
	"context"
	"tg-reminder/internal/domain"
)

type fakeUserRepo struct {
	user    *domain.User // какого пользователя вернуть
	err     error        // какую ошибку вернуть
	gotTgID int64        // сюда метод запишет, с каким tgID его вызвали
}

type fakeReminderRepo struct {
	reminders []domain.Reminder // список событий у возвращенного пользователя
	err       error             // какую ошибку вернуть
	gotUserID int64             // сюда метод запишет, с каким userID его вызвали
	gotLimit  int
	called    bool
}

func (f *fakeUserRepo) GetUserByTgID(ctx context.Context, tgID int64) (*domain.User, error) {
	f.gotTgID = tgID
	return f.user, f.err
}

func (f *fakeReminderRepo) ListUpcoming(ctx context.Context, userID int64, limit int) ([]domain.Reminder, error) {
	f.called = true
	f.gotUserID = userID
	f.gotLimit = limit
	return f.reminders, f.err
}
