package counting

import (
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/common"
	"github.com/Saime-0/tg-bot-contest/internal/model"
)

type Params struct {
	TX *sqlx.Tx

	Chat model.Chat
	User model.User
	Comp model.Contest
}

type Out struct {
	CreatedTickets []model.Ticket
}

func (p Params) Run() (Out, error) {
	// Поиск чата для которого проходит конкурс
	var competitiveChat model.Chat
	if err := p.TX.Get(&competitiveChat, `
			select chats.*
			from chats
			inner join contests 
			    on chats.id = contests.competitive_chat_id
			where contests.keyword_chat_id = ?
				and  contests.ended_at is null 
	`, p.Chat.ID); err != nil {
		return Out{}, err
	}

	var unlinkedMembers []model.Member
	if err := p.TX.Select(&unlinkedMembers, `
		select * from members
		where not ignore_in_ticket_counting 
			and in_ticket_id=0
			and status=?
			and inviter_id=?
			and chat_id=?
			and created_at>=?
	`, model.MemberStatusJoin, p.User.ID, competitiveChat.ID, p.Comp.CreatedAt,
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
		ticket := model.Ticket{
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
		memberIDs := model.MemberIDs(chunkedMembers[i])
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
