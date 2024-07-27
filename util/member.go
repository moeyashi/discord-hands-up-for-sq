package util

import "github.com/moeyashi/discord-hands-up-for-sq/repository"

func IndexOfSameMember(members []repository.Member, userID string) int {
	for index, member := range members {
		if member.UserID == userID {
			return index
		}
	}
	return -1
}
