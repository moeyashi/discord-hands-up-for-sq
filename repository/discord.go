package repository

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type discordRepository struct {
	s *discordgo.Session
}

func NewDiscordRepository(s *discordgo.Session) DiscordRepository {
	return &discordRepository{s: s}
}

func (r *discordRepository) FindRoleByName(guildID string, roleName string) (*discordgo.Role, error) {
	roles, err := r.s.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.Name == roleName ||
			//TODO 後で消す
			// https://github.com/moeyashi/discord-hands-up-for-sq/issues/34
			// での後方互換用
			// 削除用issue
			// https://github.com/moeyashi/discord-hands-up-for-sq/issues/35
			role.Name == fmt.Sprintf("内戦 %s", roleName) {
			return role, nil
		}
	}
	return nil, nil
}

func (r *discordRepository) GuildMemberRoleAdd(guildID, userID, roleID string) error {
	return r.s.GuildMemberRoleAdd(guildID, userID, roleID)
}

func (r *discordRepository) GuildMemberByRole(guildID, roleID string) ([]*discordgo.Member, error) {
	members, err := r.s.GuildMembers(guildID, "", 100)
	if err != nil {
		return nil, err
	}
	var ret []*discordgo.Member
	for _, member := range members {
		for _, role := range member.Roles {
			if role == roleID {
				ret = append(ret, member)
				break
			}
		}
	}
	return ret, nil
}
