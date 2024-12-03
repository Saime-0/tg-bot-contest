package changed

import (
	"database/sql"
	"errors"

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

	var member model.Member
	if err := p.DB.Get(&member, `
		UPDATE members
		SET status=?
		WHERE chat_id=? and user_id=?
		RETURNING *
	`, p.ParticipationType, p.Chat.ID, p.Participant.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	switch p.ParticipationType {
	case TypeJoin:
		return p.onJoin(member)
	case TypeLeave:
		return p.onLeave(member)
	default:
		return nil
	}
}

func (p *Params) onJoin(member model.Member) error {
	ticketCreatePass := p.ViaLink ||
		p.Initiator.IsBot ||
		p.Participant.IsBot ||
		p.Initiator.ID == p.Participant.ID

	if member.ID != 0 {
		return nil
	}

	member = model.Member{
		UserID:    p.Participant.ID,
		ChatID:    p.Chat.ID,
		Status:    TypeJoin,
		InviterID: 0,
	}
	if !ticketCreatePass {
		member.InviterID = p.Initiator.ID
	}

	if _, err := p.DB.NamedExec(`
			insert into members(chat_id,user_id,status,inviter_id)
			values(:chat_id,:user_id,:status,:inviter_id)
		`, member); err != nil {
		return err
	}

	return nil
}

func (p *Params) onLeave(member model.Member) error {
	if member.ID != 0 {
		return nil
	}

	member = model.Member{
		UserID:    p.Participant.ID,
		ChatID:    p.Chat.ID,
		Status:    TypeLeave,
		InviterID: 0,
	}

	if _, err := p.DB.NamedExec(`
			insert into members(chat_id,user_id,status,inviter_id)
			values(:chat_id,:user_id,:status,:inviter_id)
		`, member); err != nil {
		return err
	}

	return nil
}
