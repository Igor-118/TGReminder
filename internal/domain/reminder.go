package domain

import (
	"time"
)

type Reminder struct {
	ID       int64
	UserID   int64
	Text     string
	RemindAt time.Time
}
