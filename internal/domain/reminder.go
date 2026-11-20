package domain

import (
	"time"
)

type Reminder struct {
	ID       int64     // dateid
	UserID   int64     // user_id
	Text     string    // description
	RemindAt time.Time // time_announce
}
