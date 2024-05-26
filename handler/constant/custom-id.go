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
