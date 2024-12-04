package statusUpdate

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/model"
	chatUpdate "tgBotContest/usecase/chat/update"
	userUpdate "tgBotContest/usecase/user/update"
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
	if err := chatUpdate.Run(p.DB, p.Chat); err != nil {
		return err
	}
	if err := userUpdate.Run(p.DB, p.Participant); err != nil {
		return err
	}
	if err := userUpdate.Run(p.DB, p.Initiator); err != nil {
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
