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

func defineMemberStatus(old, new string) uint {
	oldStatus := tgModel.MemberStatusID[old]
	newStatus := tgModel.MemberStatusID[new]
	switch {
	case slices.Contains(tgModel.AlienStatus, oldStatus) && slices.Contains(tgModel.ParticipantStatus, newStatus):
		return model.MemberStatusJoin
	case slices.Contains(tgModel.ParticipantStatus, oldStatus) && slices.Contains(tgModel.AlienStatus, newStatus):
		return model.MemberStatusLeave
	default:
		return 0
	}
}

func (c *Controller) newChatMember(b *gotgbot.Bot, ctx *ext.Context) error {
	oldStatus := ctx.ChatMember.OldChatMember.GetStatus()
	newStatus := ctx.ChatMember.NewChatMember.GetStatus()
	log.Println(oldStatus, "->", newStatus)

	memberStatus := defineMemberStatus(oldStatus, newStatus)
	if memberStatus == 0 {
		return nil
	}
	initiator := ctx.ChatMember.From
	participant := ctx.ChatMember.NewChatMember.GetUser()
	viaLink := ctx.ChatMember.InviteLink != nil ||
		ctx.ChatMember.IsJoinRequest() ||
		ctx.ChatMember.ViaChatFolderInviteLink

	err := (&ucParticipationChanged.Params{
		DB:           c.DB,
		Chat:         tgModel.ChatDomain(ctx.ChatMember.Chat),
		MemberStatus: memberStatus,
		Participant:  tgModel.UserDomain(participant),
		Initiator:    tgModel.UserDomain(initiator),
		ViaLink:      viaLink,
	}).Run()
	if err != nil {
		return err
	}

	return nil
}
