package model

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
	UserID    int  `json:"user_id"`
	ChatID    int  `json:"chat_id"`
	Status    uint `json:"status"`
	InviterID int  `json:"inviter_id"`
}
