package models

import "time"

type Transactions struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Amount    int       `json:"amount"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}
