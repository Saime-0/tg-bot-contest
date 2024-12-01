package model

import (
	"github.com/PaulSonOfLars/gotgbot/v2"

	"tgBotCompetition/model"
)

func UserDomain(tgUser gotgbot.User) model.User {
	return model.User{
		ID:        int(tgUser.Id),
		FirstName: tgUser.FirstName,
		Username:  tgUser.Username,
		IsBot:     tgUser.IsBot,
	}
}
