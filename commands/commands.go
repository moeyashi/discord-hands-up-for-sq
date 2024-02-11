package commands

import "github.com/bwmarrin/discordgo"

func GetCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "husq",
			Description: "Hands up for SQ",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "list",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "SQイベントを取得します",
				},
				{
					Name:        "can",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "SQイベントに参加します",
				},
				{
					Name:        "temp",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "SQイベントに仮参加します",
				},
				{
					Name:        "sub",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "SQイベントに補欠参加します",
				},
				{
					Name:        "lounge-name",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "ラウンジでのユーザー名を表示します",
				},
				{
					Name:        "mention",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "次のSQのメンバーにメンションします",
				},
				{
					Name:        "version",
					Description: "バージョンを確認",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name: "husq set",
			Type: discordgo.MessageApplicationCommand,
		},
		// {
		// 	Name:        "civil",
		// 	Description: "Hands up for civil war",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		{
		// 			Name:        "list",
		// 			Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 			Description: "内戦イベントを取得します",
		// 		},
		// 		{
		// 			Name:        "add",
		// 			Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 			Description: "内戦イベントを追加",
		// 		},
		// 		{
		// 			Name:        "remove",
		// 			Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 			Description: "内戦イベントを削除",
		// 		},
		// 	},
		// },
		{
			Name: "setコマンドに変換",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "outコマンドに変換",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name: "sheatを保存",
			Type: discordgo.MessageApplicationCommand,
		},
	}
}
