package create

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/l10n"
	"github.com/Saime-0/tg-bot-contest/internal/model"
	"github.com/Saime-0/tg-bot-contest/internal/ue"
)

type Params struct {
	TX *sqlx.Tx

	Multiplicity      int
	Keyword           string
	CompetitiveChatID int
	KeywordChatID     int
	KeywordTopicID    int
	CreatorID         int
}

const (
	DefaultMultiplicity = 10
	DefaultKeyword      = l10n.DefaultKeyword
)

func (p *Params) Run() error {
	var competitiveChat model.Chat
	if err := p.TX.Get(&competitiveChat, `select * from chats where id = ?`, p.CompetitiveChatID); err != nil {
		return err
	}
	var keywordChat model.Chat
	if err := p.TX.Get(&keywordChat, `select * from chats where id = ?`, p.KeywordChatID); err != nil {
		return err
	}

	var exists bool
	if err := p.TX.Get(&exists, `
		select exists(
		    select 1 from contests 
			where competitive_chat_id=? 
			  and ended_at is null
		)
  	`, p.CompetitiveChatID); err != nil {
		return err
	}
	if exists {
		return ue.New(l10n.ContestCreatePreviousNotOverYet)
	}

	if p.Multiplicity <= 0 {
		p.Multiplicity = DefaultMultiplicity
	}
	if p.Keyword == "" {
		p.Keyword = DefaultKeyword
	}

	_, err := p.TX.NamedExec(`
		insert into contests (id,creator_id,competitive_chat_id,keyword_chat_id,keyword_topic_id,keyword,multiplicity)
		values (:id,:creator_id,:competitive_chat_id,:keyword_chat_id,:keyword_topic_id,:keyword,:multiplicity)
	`, model.Contest{
		ID:                uuid.NewString(),
		CreatorID:         p.CreatorID,
		CompetitiveChatID: p.CompetitiveChatID,
		KeywordChatID:     p.KeywordChatID,
		KeywordTopicID:    p.KeywordTopicID,
		Keyword:           p.Keyword,
		Multiplicity:      p.Multiplicity,
		//CreatedAt:         now,
		//EndedAt:           nil,
	})

	return err
}
