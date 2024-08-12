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
		{
			Name:        "mogi",
			Description: "Hands up for civil war",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "list",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "内戦イベントを取得します",
				},
				{
					Name:        "set",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "内戦イベントを追加",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "month",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Description: "月",
							Required:    true,
							MinValue:    makeFloat64(1),
							MaxValue:    float64(12),
						},
						{
							Name:        "date",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Description: "日",
							Required:    true,
							MinValue:    makeFloat64(1),
							MaxValue:    float64(31),
						},
						{
							Name:        "hour",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Description: "時",
							Required:    true,
							MinValue:    makeFloat64(1),
							MaxValue:    float64(24),
						},
					},
				},
				{
					Name:        "remove",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "内戦イベントを削除",
				},
				{
					Name:        "can",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "内戦イベントに参加します",
				},
				{
					Name:        "temp",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "内戦イベントに仮参加します",
				},
				{
					Name:        "sub",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "内戦イベントに補欠参加します",
				},
			},
		},
		{
			Name: "sheatを保存",
			Type: discordgo.MessageApplicationCommand,
		},
		{
			Name:        "results",
			Description: "戦績の管理",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "url",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "スプレッドシートのURLを取得します",
				},
				{
					Name:        "set-url",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "スプレッドシートのURLを設定します",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "url",
							Type:        discordgo.ApplicationCommandOptionString,
							Description: "スプレッドシートのURL",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:        "mkmg",
			Description: "MKMG",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "team",
					Description: "チーム名を設定します",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "name",
							Description: "チーム名",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    false,
						},
					},
				},
				{
					Name:        "post",
					Description: "mkmgの投稿文を作成します",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "time",
							Description: "時間（例：21）",
							Type:        discordgo.ApplicationCommandOptionString,
							Required:    true,
						},
						{
							Name:        "avg",
							Description: "平均（例：5000）",
							Type:        discordgo.ApplicationCommandOptionInteger,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func makeFloat64(i int) *float64 {
	f := float64(i)
	return &f
}
