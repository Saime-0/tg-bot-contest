package model

import (
	"github.com/PaulSonOfLars/gotgbot/v2"

	"github.com/Saime-0/tg-bot-contest/internal/model"
)

func UserDomain(tgUser gotgbot.User) model.User {
	return model.User{
		ID:        int(tgUser.Id),
		FirstName: tgUser.FirstName,
		Username:  tgUser.Username,
		IsBot:     tgUser.IsBot,
	}
}

func ChatFullDomain(tgChat *gotgbot.Chat) model.Chat {
	return model.Chat{
		ID:       int(tgChat.Id),
		Title:    tgChat.Title,
		Username: tgChat.Username,
	}
}
