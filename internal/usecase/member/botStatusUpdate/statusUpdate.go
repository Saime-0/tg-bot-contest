package botStatusUpdate

import (
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/Saime-0/tg-bot-contest/internal/model"
	chatUpdate "github.com/Saime-0/tg-bot-contest/internal/usecase/chat/update"
	contestStop "github.com/Saime-0/tg-bot-contest/internal/usecase/contests/stop"
)

type Params struct {
	TX *sqlx.Tx

	Chat            model.Chat
	LinkedChatID    int
	BotMemberStatus uint
}

func (p *Params) Run() error {
	if err := chatUpdate.Run(p.TX, p.Chat); err != nil {
		return err
	}

	// Если бот покинул чат, остановить текущий конкурс
	if p.BotMemberStatus == model.MemberStatusLeave {
		err := (&contestStop.Params{
			TX:     p.TX,
			ChatID: p.Chat.ID,
		}).Run()
		// Если ошибки нет либо конкурс не запущен, посчитаем что ошибки нет
		if err == nil || errors.Is(err, contestStop.ErrContestNotFound) {
			return nil
		}
		return err
	}

	return nil
}
