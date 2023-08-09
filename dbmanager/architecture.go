package dbmanager

import "time"

type Expense struct {
	ID         int64
	CategoryID int64
	Amount     float32
	Date       time.Time
	Location   string
	Notes      string
}

type Category struct {
	ID          int64
	Name        string
	Description string
}

type Tag struct {
	ID   int64
	Name string
}

type ExpenseTag struct {
	ID        int64
	ExpenseID int64
	TagId     int64
}
