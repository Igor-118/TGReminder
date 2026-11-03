package service_test

import (
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
	gotUserID int64             // сюда метод запишет, с каким tgID его вызвали
	gotLimit  int
	called    bool
}
