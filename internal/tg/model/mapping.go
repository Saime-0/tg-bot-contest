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

func ChatFullDomain(tgChat *gotgbot.ChatFullInfo) model.Chat {
	chat := model.Chat{
		ID:       int(tgChat.Id),
		Title:    tgChat.Title,
		Username: tgChat.Username,
		ChildID:  0,
		ParentID: 0,
	}

	switch tgChat.Type {
	// Если этот чат является группой, то связанный чат является родительским
	case gotgbot.ChatTypeGroup,
		gotgbot.ChatTypeSupergroup:
		chat.ParentID = int(tgChat.LinkedChatId)
	// Если этот чат является каналом, то связанный чат является дочерним
	case gotgbot.ChatTypeChannel:
		chat.ChildID = int(tgChat.LinkedChatId)
	// В Остальных случаях ничего не происходит
	case gotgbot.ChatTypePrivate:
	}

	return chat
}
