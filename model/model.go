package model

import "time"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

const (
	MemberStatusJoin  uint = 1
	MemberStatusLeave uint = 2
)

type Member struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	ChatID    int  `json:"chat_id"`
	Status    uint `json:"status"`
	InviterID int  `json:"inviter_id"`
}

type Chat struct {
	ID       int    `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	Username string `json:"username" db:"username"`
}

type Competition struct {
	ID        int        `json:"id" db:"id"`
	CreatorID int        `json:"creator_id" db:"creator_id"`
	ChatID    int        `json:"chat_id" db:"chat_id"`
	TopicID   int        `json:"topic_id" db:"topic_id"`
	Keyword   string     `json:"keyword" db:"keyword"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	EndedAt   *time.Time `json:"ended_at" db:"ended_at"`
}
