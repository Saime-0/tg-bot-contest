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

	Chat         model.Chat
	MemberStatus uint
	Participant  model.User
	Initiator    model.User
	ViaLink      bool
}

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
	`, p.MemberStatus, p.Chat.ID, p.Participant.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if member.ID == 0 {
		return p.saveMember(member)
	}

	return nil
}

func (p *Params) saveMember(member model.Member) error {
	ignoreInTicketCounting := p.MemberStatus == model.MemberStatusLeave ||
		p.ViaLink ||
		p.Initiator.IsBot ||
		p.Participant.IsBot ||
		p.Initiator.ID == p.Participant.ID

	member = model.Member{
		UserID:                 p.Participant.ID,
		ChatID:                 p.Chat.ID,
		Status:                 p.MemberStatus,
		InviterID:              0,
		IgnoreInTicketCounting: ignoreInTicketCounting,
	}

	if p.MemberStatus == model.MemberStatusJoin {
		member.InviterID = p.Initiator.ID
	}

	if _, err := p.DB.NamedExec(`
			insert into members(chat_id,user_id,status,inviter_id,ignore_in_ticket_counting)
			values(:chat_id,:user_id,:status,:inviter_id,:ignore_in_ticket_counting)
		`, member); err != nil {
		return err
	}

	return nil
}
