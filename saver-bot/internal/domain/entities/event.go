package entities

type Manic struct {
	ID   string `db:"id"`
	Date string `db:"event_date"`
}

type Massage struct {
	ID   string `db:"id"`
	Date string `db:"event_date"`
}
