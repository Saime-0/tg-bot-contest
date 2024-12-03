package changed

import (
	"database/sql"
	"errors"
	"slices"

	"github.com/jmoiron/sqlx"

	"tgBotCompetition/model"
	chatUpdate "tgBotCompetition/usecase/chat/update"
	userUpdate "tgBotCompetition/usecase/user/update"
)

type Params struct {
	DB *sqlx.DB

	Chat              model.Chat
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
	if err := (&chatUpdate.Params{DB: p.DB, Chat: p.Chat}).Run(); err != nil {
		return err
	}
	if err := (&userUpdate.Params{DB: p.DB, User: p.Participant}).Run(); err != nil {
		return err
	}
	if err := (&userUpdate.Params{DB: p.DB, User: p.Initiator}).Run(); err != nil {
		return err
	}

	if p.ViaLink ||
		p.Initiator.IsBot ||
		p.Participant.IsBot ||
		p.Initiator.ID == p.Participant.ID {
		return nil
	}

	var member model.Member
	if err := p.DB.Get(&member, `
		select * from members
		where chat_id=? and user_id=?
	`, p.Chat.ID, p.Participant.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if member.ID != 0 && slices.Contains()member.Status ==  {
		return nil
	}

	return nil
}
