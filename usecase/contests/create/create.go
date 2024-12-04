package create

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"tgBotContest/l10n"
	"tgBotContest/model"
)

type Params struct {
	DB *sqlx.DB

	Multiplicity int
	Keyword      string
	ChatID       int
	TopicID      int
	CreatorID    int
}

const (
	DefaultMultiplicity = 10
	DefaultKeyword      = l10n.DefaultKeyword
)

var (
	ErrContestAlreadyExists = errors.New("конкурс уже запущен")
)

func (p *Params) Run() error {
	var exists bool
	if err := p.DB.Get(&exists, `
		select exists(
		    select 1 from contests 
			where chat_id=? 
			  and ended_at is null
		)
  	`, p.ChatID); err != nil {
		return err
	}
	if exists {
		return ErrContestAlreadyExists
	}

	if p.Multiplicity <= 0 {
		p.Multiplicity = DefaultMultiplicity
	}
	if p.Keyword == "" {
		p.Keyword = DefaultKeyword
	}

	_, err := p.DB.NamedExec(`
		insert into contests (creator_id,chat_id,topic_id,keyword,multiplicity)
		values (:creator_id,:chat_id,:topic_id,:keyword,:multiplicity)`,
		model.Contest{
			CreatorID:    p.CreatorID,
			ChatID:       p.ChatID,
			TopicID:      p.TopicID,
			Keyword:      p.Keyword,
			Multiplicity: p.Multiplicity,
		})

	return err
}
