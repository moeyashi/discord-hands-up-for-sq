package usecase

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/domain/discord"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
	"github.com/moeyashi/discord-hands-up-for-sq/util"
)

func AppendMogiMember(ctx context.Context, repo repository.Repository, discordRepo repository.DiscordRepository, guild *repository.Guild, mogiTitle string, member *discordgo.Member, memberType repository.MemberTypes) error {
	mogi, err := repo.GetMogi(ctx, guild, mogiTitle)
	if err != nil {
		return err
	}

	role, err := discordRepo.FindRoleByName(guild.ID, mogi.RoleName())
	if err != nil {
		return err
	}
	if role != nil {
		if err := discordRepo.GuildMemberRoleAdd(guild.ID, member.User.ID, role.ID); err != nil {
			return err
		}
	}

	members := mogi.Members
	existsIndex := util.IndexOfSameMember(members, member.User.ID)
	if existsIndex >= 0 {
		// 既に参加している場合は一旦削除
		members = append(members[:existsIndex], members[existsIndex+1:]...)
	}
	members = append(members, repository.Member{
		UserID:     member.User.ID,
		UserName:   discord.GetDisplayUsername(member),
		MemberType: memberType,
	})
	if err := repo.PutMogiMembers(ctx, guild, mogiTitle, members); err != nil {
		return err
	}

	return nil
}
