package model

type Birthday struct {
	Day   int64  `db:"day"`
	Month int64  `db:"month"`
	Year  *int64 `db:"year"`
}
