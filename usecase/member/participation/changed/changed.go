package changed

import (
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	ChatID            int
	ParticipationType int
	Participant       model.User
	Initiator         model.User
	ViaLink           bool
}

const (
	TypeJoin  = 1
	TypeLeave = 2
)

func (p *Params) Run() error {
	if p.ViaLink ||
		p.Initiator.IsBot ||
		p.Participant.IsBot ||
		p.Initiator.ID == p.Participant.ID {
		return nil
	}

	return nil
}
