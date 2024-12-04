package statusUpdate

import (
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

	if res, err := p.DB.Exec(`
		update members
		set status=?
		where chat_id=? and user_id=?
	`, p.MemberStatus, p.Chat.ID, p.Participant.ID); err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected == 0 {
		return p.saveMember()
	}

	return nil
}

func (p *Params) saveMember() error {
	ignoreInTicketCounting := p.MemberStatus == model.MemberStatusLeave ||
		p.ViaLink ||
		p.Initiator.IsBot ||
		p.Participant.IsBot ||
		p.Initiator.ID == p.Participant.ID

	member := model.Member{
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
