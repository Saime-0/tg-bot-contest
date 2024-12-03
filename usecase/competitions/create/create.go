package create

import (
	"errors"

	"github.com/jmoiron/sqlx"

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
	_, err := p.DB.NamedExec(`
		INSERT INTO competitions (creator_id,chat_id,topic_id)
		VALUES (:creator_id,:chat_id,:topic_id)`,
		model.Competition{
			CreatorID: p.CreatorID,
			ChatID:    p.ChatID,
			TopicID:   p.TopicID,
			Keyword:   p.Keyword,
		})

	return err
}
