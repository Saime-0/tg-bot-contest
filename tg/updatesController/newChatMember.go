package updatesController

import (
	"log"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"

	"tgBotCompetition/model"
	tgModel "tgBotCompetition/tg/model"
	ucParticipationChanged "tgBotCompetition/usecase/member/participation/changed"
)

var isMemberStatus = []int{
	model.MemberStatusMember,
	model.MemberStatusRestricted,
	model.MemberStatusAdministrator,
	model.MemberStatusCreator,
}

var isOutStatus = []int{
	model.MemberStatusLeft,
	model.MemberStatusKicked,
}

func defineParticipationType(old, new string) int {
	oldStatus := model.MemberStatusID[old]
	newStatus := model.MemberStatusID[new]
	switch {
	case slices.Contains(isOutStatus, oldStatus) && slices.Contains(isMemberStatus, newStatus):
		return ucParticipationChanged.TypeJoin
	case slices.Contains(isMemberStatus, oldStatus) && slices.Contains(isOutStatus, newStatus):
		return ucParticipationChanged.TypeLeave
	default:
		return 0
	}
}

func (c *Controller) newChatMember(b *gotgbot.Bot, ctx *ext.Context) error {
	oldStatus := ctx.ChatMember.OldChatMember.GetStatus()
	newStatus := ctx.ChatMember.NewChatMember.GetStatus()
	log.Println(oldStatus, "->", newStatus)

	participationType := defineParticipationType(oldStatus, newStatus)
	if participationType == 0 {
		return nil
	}
	initiator := ctx.ChatMember.From
	participant := ctx.ChatMember.NewChatMember.GetUser()
	viaLink := ctx.ChatMember.InviteLink != nil ||
		ctx.ChatMember.IsJoinRequest() ||
		ctx.ChatMember.ViaChatFolderInviteLink

	err := (&ucParticipationChanged.Params{
		DB:                c.DB,
		ChatID:            int(ctx.ChatMember.Chat.Id),
		ParticipationType: participationType,
		Participant:       tgModel.UserDomain(participant),
		Initiator:         tgModel.UserDomain(initiator),
		ViaLink:           viaLink,
	}).Run()
	if err != nil {
		return err
	}

	return nil
}
