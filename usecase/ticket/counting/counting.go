package counting

import (
	"github.com/jmoiron/sqlx"

	"tgBotCompetition/common"
	"tgBotCompetition/model"
)

type Params struct {
	DB *sqlx.DB

	Chat model.Chat
	User model.User
	Comp model.Competition
}

type Out struct {
	CreatedTickets []model.Ticket
}

func (p Params) Run() (Out, error) {
	var unlinkedMembers []model.Member
	if err := p.DB.Select(&unlinkedMembers, `
		select * from members
		where not ignore_in_ticket_counting 
			and in_ticket_id=0
			and status=?
			and inviter_id=?
			and chat_id=?
			and created_at>=?
	`, model.MemberStatusJoin, p.User.ID, p.Chat.ID, p.Comp.CreatedAt,
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
	if err := p.DB.Get(&lastTicketNumber, `
		select ifnull(max(number),0) from tickets 
		where competition_id=?
	`, p.Comp.ID); err != nil {
		return Out{}, err
	}

	var out Out
	for i := 0; i < len(chunkedMembers); i++ {
		ticket := model.Ticket{
			Number:        i + lastTicketNumber + 1,
			UserID:        p.User.ID,
			CompetitionID: p.Comp.ID,
		}
		if _, err := p.DB.NamedExec(`
			insert into tickets(number, user_id, competition_id)
			values (:number, :user_id, :competition_id)
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
		} else if _, err = p.DB.Exec(q, args...); err != nil {
			return Out{}, err
		} else {
			out.CreatedTickets = append(out.CreatedTickets, ticket)
		}
	}

	return out, nil
}
