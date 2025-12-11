package domain

import "time"

type Memo struct {
	ID int64 `json:"id"`
	UserID string `json:"user_id"`
	Body string `json:"body"`
	Mood string `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}