package create

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"tgBotCompetition/l10n"
	"tgBotCompetition/model"
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
	ErrCompetitionAlreadyExists = errors.New("конкурс уже запущен")
)

func (p *Params) Run() error {
	var exists bool
	if err := p.DB.Select(&exists, `
		SELECT 1
		FROM competitions 
		WHERE chat_id = ?
		  AND ended_at IS NOT NULL`); err != nil {
		return err
	}
	if exists {
		return ErrCompetitionAlreadyExists
	}

	if p.Multiplicity <= 0 {
		p.Multiplicity = DefaultMultiplicity
	}
	if p.Keyword == "" {
		p.Keyword = DefaultKeyword
	}

	_, err := p.DB.NamedExec(`
		INSERT INTO competitions (creator_id,chat_id,topic_id,keyword,multiplicity)
		VALUES (:creator_id,:chat_id,:topic_id,:keyword,:multiplicity)`,
		model.Competition{
			CreatorID:    p.CreatorID,
			ChatID:       p.ChatID,
			TopicID:      p.TopicID,
			Keyword:      p.Keyword,
			Multiplicity: p.Multiplicity,
		})

	return err
}
