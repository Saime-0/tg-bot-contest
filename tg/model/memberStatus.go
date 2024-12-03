package model

const (
	MemberStatusLeft          = 1
	MemberStatusKicked        = 2
	MemberStatusMember        = 3
	MemberStatusRestricted    = 4
	MemberStatusAdministrator = 5
	MemberStatusCreator       = 6
)

var MemberStatusNameMap = map[uint]string{
	MemberStatusLeft:          "left",
	MemberStatusKicked:        "kicked",
	MemberStatusMember:        "member",
	MemberStatusRestricted:    "restricted",
	MemberStatusAdministrator: "administrator",
	MemberStatusCreator:       "creator",
}

var MemberStatusID = map[string]uint{
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

var ParticipantStatus = []uint{
	MemberStatusMember,
	MemberStatusRestricted,
	MemberStatusAdministrator,
	MemberStatusCreator,
}

var AlienStatus = []uint{
	MemberStatusLeft,
	MemberStatusKicked,
}
