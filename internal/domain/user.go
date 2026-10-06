package domain

type User struct {
	ID       int64  `db:"id"`
	TgID     int64  `db:"tg_id"`
	Name     string `db:"name"`
	Age      int    `db:"age"`
	TimeZone string `db:"time_zone"`
}
