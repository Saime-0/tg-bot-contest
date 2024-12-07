package counting

import (
	"github.com/jmoiron/sqlx"

	"tgBotContest/internal/common"
	model2 "tgBotContest/internal/model"
)

type Params struct {
	TX *sqlx.Tx

	Chat model2.Chat
	User model2.User
	Comp model2.Contest
}

type Out struct {
	CreatedTickets []model2.Ticket
}

func (p Params) Run() (Out, error) {
	var unlinkedMembers []model2.Member
	if err := p.TX.Select(&unlinkedMembers, `
		select * from members
		where not ignore_in_ticket_counting 
			and in_ticket_id=0
			and status=?
			and inviter_id=?
			and chat_id=?
			and created_at>=?
	`, model2.MemberStatusJoin, p.User.ID, p.Chat.ID, p.Comp.CreatedAt,
	); err != nil {
		return Out{}, err
	}
	if len(unlinkedMembers) == 0 {
		return Out{}, nil
	}

	if len(unlinkedMembers)/p.Comp.Multiplicity == 0 {
		return Out{}, nil
	}

	chunkedMembers := common.ChunkSlice(unlinkedMembers, p.Comp.Multiplicity)

	var lastTicketNumber int
	if err := p.TX.Get(&lastTicketNumber, `
		select ifnull(max(number),0) from tickets 
		where contest_id=?
	`, p.Comp.ID); err != nil {
		return Out{}, err
	}

	var out Out
	for i := 0; i < len(chunkedMembers); i++ {
		ticket := model2.Ticket{
			Number:    i + lastTicketNumber + 1,
			UserID:    p.User.ID,
			ContestID: p.Comp.ID,
		}
		if _, err := p.TX.NamedExec(`
			insert into tickets(number, user_id, contest_id)
			values (:number, :user_id, :contest_id)
		`, ticket); err != nil {
			return Out{}, err
		}
		memberIDs := model2.MemberIDs(chunkedMembers[i])
		if q, args, err := sqlx.In(`
			update members
			set in_ticket_id=?
			where id in (?)
		`, ticket.Number, memberIDs); err != nil {
			return Out{}, err
		} else if _, err = p.TX.Exec(q, args...); err != nil {
			return Out{}, err
		} else {
			out.CreatedTickets = append(out.CreatedTickets, ticket)
		}
	}

	return out, nil
}
