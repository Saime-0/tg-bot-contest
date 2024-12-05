package updatesController

import (
	"log"
	"slices"

	"tgBotContest/model"
	tgModel "tgBotContest/tg/model"
	memberStatusUpdate "tgBotContest/usecase/member/statusUpdate"
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

func newChatMember(r Request) error {
	oldStatus := r.ctx.ChatMember.OldChatMember.GetStatus()
	newStatus := r.ctx.ChatMember.NewChatMember.GetStatus()
	log.Println(oldStatus, "->", newStatus)

	memberStatus := defineMemberStatus(oldStatus, newStatus)
	if memberStatus == 0 {
		return nil
	}
	initiator := r.ctx.ChatMember.From
	participant := r.ctx.ChatMember.NewChatMember.GetUser()
	viaLink := r.ctx.ChatMember.InviteLink != nil ||
		r.ctx.ChatMember.IsJoinRequest() ||
		r.ctx.ChatMember.ViaChatFolderInviteLink

	err := (&memberStatusUpdate.Params{
		TX:           r.TX(),
		Chat:         tgModel.ChatDomain(r.ctx.ChatMember.Chat),
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
