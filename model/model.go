package model

import "time"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

const (
	MemberStatusLeft          = 1
	MemberStatusKicked        = 2
	MemberStatusMember        = 3
	MemberStatusRestricted    = 4
	MemberStatusAdministrator = 5
	MemberStatusCreator       = 6
)

var MemberStatusNameMap = map[int]string{
	MemberStatusLeft:          "left",
	MemberStatusKicked:        "kicked",
	MemberStatusMember:        "member",
	MemberStatusRestricted:    "restricted",
	MemberStatusAdministrator: "administrator",
	MemberStatusCreator:       "creator",
}

var MemberStatusID = map[string]int{
	"left":          MemberStatusLeft,
	"kicked":        MemberStatusKicked,
	"member":        MemberStatusMember,
	"restricted":    MemberStatusRestricted,
	"administrator": MemberStatusAdministrator,
	"creator":       MemberStatusCreator,
}

func MemberStatusName(memberID int) string {
	switch memberID {
	case MemberStatusLeft:
		return "left"
	case MemberStatusKicked:
		return "kicked"
	case MemberStatusMember:
		return "member"
	case MemberStatusRestricted:
		return "restricted"
	case MemberStatusAdministrator:
		return "administrator"
	case MemberStatusCreator:
		return "creator"
	default:
		return "unknown"
	}

}

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

var ParticipantStatus = []int{
	MemberStatusMember,
	MemberStatusRestricted,
	MemberStatusAdministrator,
	MemberStatusCreator,
}

var AlienStatus = []int{
	MemberStatusLeft,
	MemberStatusKicked,
}
