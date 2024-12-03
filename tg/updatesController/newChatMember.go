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

func defineParticipationType(old, new string) int {
	oldStatus := model.MemberStatusID[old]
	newStatus := model.MemberStatusID[new]
	switch {
	case slices.Contains(model.AlienStatus, oldStatus) && slices.Contains(model.ParticipantStatus, newStatus):
		return ucParticipationChanged.TypeJoin
	case slices.Contains(model.ParticipantStatus, oldStatus) && slices.Contains(model.AlienStatus, newStatus):
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
		Chat:              tgModel.ChatDomain(ctx.ChatMember.Chat),
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
