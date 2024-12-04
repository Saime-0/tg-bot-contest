package model

import (
	"github.com/PaulSonOfLars/gotgbot/v2"

	"tgBotContest/model"
)

func UserDomain(tgUser gotgbot.User) model.User {
	return model.User{
		ID:        int(tgUser.Id),
		FirstName: tgUser.FirstName,
		Username:  tgUser.Username,
		IsBot:     tgUser.IsBot,
	}
}

func ChatDomain(tgChat gotgbot.Chat) model.Chat {
	return model.Chat{
		ID:       int(tgChat.Id),
		Title:    tgChat.FirstName,
		Username: tgChat.Username,
	}
}
