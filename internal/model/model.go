package model

import "time"

type User struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	Username  string `db:"username"`
	IsBot     bool   `db:"is_bot"`
}

const (
	MemberStatusJoin  uint = 1
	MemberStatusLeave uint = 2
)

type Member struct {
	ID                     int  `db:"id"`
	UserID                 int  `db:"user_id"`
	ChatID                 int  `db:"chat_id"`
	Status                 uint `db:"status"`
	InviterID              int  `db:"inviter_id"`
	IgnoreInTicketCounting bool `db:"ignore_in_ticket_counting"`
	InTicketID             int  `db:"in_ticket_id"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Chat struct {
	ID       int    `db:"id"`
	Title    string `db:"title"`
	Username string `db:"username"`

	CreatedAt time.Time `db:"created_at"`
}

type Contest struct {
	ID                string     `db:"id"`
	CreatorID         int        `db:"creator_id"`
	CompetitiveChatID int        `db:"competitive_chat_id"`
	KeywordChatID     int        `db:"keyword_chat_id"`
	KeywordTopicID    int        `db:"keyword_topic_id"`
	Keyword           string     `db:"keyword"`
	Multiplicity      int        `db:"multiplicity"`
	CreatedAt         time.Time  `db:"created_at"`
	EndedAt           *time.Time `db:"ended_at"`
}

type Ticket struct {
	Number    int    `json:"number" db:"number"`
	UserID    int    `json:"user_id" db:"user_id"`
	ContestID string `json:"contest_id" db:"contest_id"`
}
