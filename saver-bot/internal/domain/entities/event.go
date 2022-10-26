package entities

type Event struct {
	ID       string `db:"id"`
	Date     string `db:"date"`
	Type     string `db:"type"`
	TelegaID int64  `db:"telega_id"`
}

type EventAll struct {
	Date string
	Type string
}
