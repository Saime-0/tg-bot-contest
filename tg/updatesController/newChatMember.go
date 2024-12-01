package updatesController

import (
	"log"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	ucJoin "tgBotCompetition/usecase/member/join"
	ucLeft "tgBotCompetition/usecase/member/left"
)

const (
	left          = "left"
	kicked        = "kicked"
	member        = "member"
	restricted    = "restricted"
	administrator = "administrator"
	creator       = "creator"
)

const (
	Unknown = iota
	IsJoinAction
	IsLeaveAction
)

var isMemberStatus = []string{member, restricted, administrator, creator}
var isOutStatus = []string{left, kicked}

func defineAction(old, new string) int {
	switch {
	case slices.Contains(isOutStatus, old) && slices.Contains(isMemberStatus, new):
		return IsJoinAction
	case slices.Contains(isMemberStatus, old) && slices.Contains(isOutStatus, new):
		return IsLeaveAction
	default:
		return Unknown
	}
}

func (c *Controller) newChatMember(b *gotgbot.Bot, ctx *ext.Context) error {
	oldStatus := ctx.ChatMember.OldChatMember.GetStatus()
	newStatus := ctx.ChatMember.NewChatMember.GetStatus()
	log.Println(oldStatus, "->", newStatus)

	switch defineAction(oldStatus, newStatus) {
	case IsJoinAction:
		return (&ucJoin.Params{}).Run()
	case IsLeaveAction:
		return (&ucLeft.Params{}).Run()
	default:
		return nil
	}
}
