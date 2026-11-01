package domain

import (
	"time"
)

type Reminder struct {
	ID       int64
	TgID     int64
	Text     string
	RemindAt time.Time
}
