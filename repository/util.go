package repository

import (
	"os"

	"google.golang.org/api/option"
)

// UserID と MemberType が一致する Member のインデックスを返す
//
// 一致する Member が存在しない場合は -1 を返す
func IndexOfSameRegistered(members []Member, userID string, memberType MemberTypes) int {
	for index, member := range members {
		if member.UserID == userID && member.MemberType == memberType {
			return index
		}
	}
	return -1
}

var credentialJSON string

func init() {
	credentialJSON = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
}

func GetGoogleDefaultCredentialClientOption() option.ClientOption {
	return option.WithCredentialsJSON([]byte(credentialJSON))
}
