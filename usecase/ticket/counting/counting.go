package counting

import (
	"github.com/jmoiron/sqlx"

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
	if err := p.DB.Get(&unlinkedMembers, `
		select * from members
		where not ignore_in_ticket_counting 
			and in_ticket_id is null
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

	ticketsNumber := len(unlinkedMembers) / p.Comp.Multiplicity
	if ticketsNumber == 0 {
		return nil
	}

	var lastTicketNumber int
	if err := p.DB.Get(&lastTicketNumber, `
		select max(number) from tickets 
		where competition_id=?
	`); err != nil {
		return err
	}

	tickets := make([]model.Ticket, ticketsNumber)
	for i := range tickets {
		tickets[i] = model.Ticket{
			Number:        i + lastTicketNumber + 1,
			UserID:        p.User.ID,
			CompetitionID: p.Comp.ID,
		}
	}
	if _, err := p.DB.NamedExec(`
		insert into tickets(number, user_id, competition_id)
		values (:number, :user_id, :competition_id)
	`, tickets); err != nil {
		return err
	}

	return nil
}
