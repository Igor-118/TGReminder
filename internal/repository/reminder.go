package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tg-reminder/internal/domain"
)

type ReminderRepository struct {
	pool *pgxpool.Pool
}

func NewReminderRepository(pool *pgxpool.Pool) *ReminderRepository {
	return &ReminderRepository{pool: pool}
}

func (r *ReminderRepository) ListUpcoming(ctx context.Context, userID int64,
	limit int) ([]domain.Reminder, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT dateid, user_id, description, time_announce
        FROM reminders
        WHERE user_id = $1
          AND time_announce > NOW()
        ORDER BY time_announce ASC
        LIMIT $2
    `, userID, limit)
	if err != nil {
		return nil, fmt.Errorf("list upcoming reminders: %w", err)
	}
	defer rows.Close()
	reminders := make([]domain.Reminder, 0, limit)
	for rows.Next() {
		var reminder domain.Reminder
		err := rows.Scan(
			&reminder.ID,
			&reminder.UserID,
			&reminder.Text,
			&reminder.RemindAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan reminder: %w", err)
		}
		reminders = append(reminders, reminder)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reminders: %w", err)
	}

	return reminders, nil
}
