package counting

import (
	"github.com/jmoiron/sqlx"
	"github.com/nullism/bqb"

	"tgBotCompetition/common"
	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	Chat model.Chat
	User model.User
	Comp model.Competition
}

func (p Params) Run() error {
	var unlinkedMembers []model.Member
	if err := p.DB.Select(&unlinkedMembers, `
		select * from members
		where not ignore_in_ticket_counting 
			and in_ticket_id=0
			and status=?
			and inviter_id=?
			and chat_id=?
	`, model.MemberStatusJoin, p.User.ID, p.Chat.ID,
	); err != nil {
		return err
	}
	if len(unlinkedMembers) == 0 {
		return nil
	}

	if len(unlinkedMembers)/p.Comp.Multiplicity == 0 {
		return nil
	}

	chunkedMembers := common.ChunkSlice(unlinkedMembers, p.Comp.Multiplicity)

	var lastTicketNumber int
	if err := p.DB.Get(&lastTicketNumber, `
		select ifnull(max(number),0) from tickets 
		where competition_id=?
	`, p.Comp.ID); err != nil {
		return err
	}

	for i := 0; i < len(chunkedMembers); i++ {
		ticketNumber := i + lastTicketNumber + 1
		if _, err := p.DB.Exec(`
			insert into tickets(number, user_id, competition_id)
			values (?, ?, ?)
		`, ticketNumber, p.User.ID, p.Comp.ID); err != nil {
			return err
		}
		memberIDs := model.MemberIDs(chunkedMembers[i])
		if q, err := bqb.New(`
			update members
			set in_ticket_id=?
			where id in (?)
		`, ticketNumber, memberIDs).ToRaw(); err != nil {
			return err
		} else if _, err = p.DB.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

//bqb.New()
//memberIDs := model.MemberIDs(chunkedMembers[i])
//if q, args, err := sqlx.In(`
//			update members
//			set in_ticket_id=?
//			where id in (?)
//		`, ticketNumber, memberIDs); err != nil {
//return err
//} else if _, err = p.DB.Query(q, args); err != nil {
//return err
//}
