package constant

import "github.com/moeyashi/discord-hands-up-for-sq/repository"

type SQListSelectCustomID string

const (
	SQListSelectCustomIDCan  SQListSelectCustomID = "can_select"
	SQListSelectCustomIDTemp SQListSelectCustomID = "temp_select"
	SQListSelectCustomIDSub  SQListSelectCustomID = "sub_select"
)

func SQListSelectCustomIDFromString(customID string) SQListSelectCustomID {
	switch customID {
	case string(SQListSelectCustomIDCan):
		return SQListSelectCustomIDCan
	case string(SQListSelectCustomIDTemp):
		return SQListSelectCustomIDTemp
	case string(SQListSelectCustomIDSub):
		return SQListSelectCustomIDSub
	default:
		return SQListSelectCustomIDCan
	}
}

func (c SQListSelectCustomID) ToMemberTypes() repository.MemberTypes {
	switch c {
	case SQListSelectCustomIDCan:
		return repository.MemberTypesParticipant
	case SQListSelectCustomIDTemp:
		return repository.MemberTypesTemporary
	case SQListSelectCustomIDSub:
		return repository.MemberTypesSub
	default:
		return repository.MemberTypesParticipant
	}
}

type MogiListSelectCustomID string

const (
	MogiListSelectCustomIDCan  MogiListSelectCustomID = "mogi_select_can"
	MogiListSelectCustomIDTemp MogiListSelectCustomID = "mogi_select_temp"
	MogiListSelectCustomIDSub  MogiListSelectCustomID = "mogi_select_sub"
)

func MogiListSelectCustomIDFromString(customID string) MogiListSelectCustomID {
	switch customID {
	case string(MogiListSelectCustomIDCan):
		return MogiListSelectCustomIDCan
	case string(MogiListSelectCustomIDTemp):
		return MogiListSelectCustomIDTemp
	case string(MogiListSelectCustomIDSub):
		return MogiListSelectCustomIDSub
	default:
		return MogiListSelectCustomIDCan
	}
}

func (c MogiListSelectCustomID) ToMemberTypes() repository.MemberTypes {
	switch c {
	case MogiListSelectCustomIDCan:
		return repository.MemberTypesParticipant
	case MogiListSelectCustomIDTemp:
		return repository.MemberTypesTemporary
	case MogiListSelectCustomIDSub:
		return repository.MemberTypesSub
	default:
		return repository.MemberTypesParticipant
	}
}
