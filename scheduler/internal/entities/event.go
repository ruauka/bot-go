package entities

type Event struct {
	ID             string `db:"id"`
	Date           string `db:"date"`
	Whom           string `db:"whom"`
	Type           string `db:"type"`
	Username       string `db:"username"`
	TelegaID       int64  `db:"telega_id"`
	ReminderStatus int
}
