package main

import (
	"context"
	"flag"
	"strings"

	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/commands"
	"github.com/moeyashi/discord-hands-up-for-sq/handler"
	"github.com/moeyashi/discord-hands-up-for-sq/handler/constant"
	"github.com/moeyashi/discord-hands-up-for-sq/repository"
)

// 参考：https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

func init() {
	flag.Parse()

	EnvBotToken := os.Getenv("BOT_TOKEN")
	if EnvBotToken != "" {
		*BotToken = EnvBotToken
	}
}

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

var (
	commandHandlers = map[string]func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository){
		"husq set": handler.SetSQ,
		"husq": func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
			options := i.ApplicationCommandData().Options
			switch options[0].Name {
			case "list":
				handler.ListSQ(ctx, s, i, repository)
				return
			case "can":
				handler.HandleCan(ctx, s, i, repository)
				return
			case "temp":
				handler.HandleTemp(ctx, s, i, repository)
				return
			case "sub":
				handler.HandleSub(ctx, s, i, repository)
				return
			case "lounge-name":
				handler.HandleLoungeName(ctx, s, i, repository)
				return
			case "mention":
				handler.HandleMention(ctx, s, i, repository)
				return
			case "version":
				handler.GetVersion(ctx, s, i, repository)
				return
			}
		},
		"setコマンドに変換": handler.CreateSetCommands,
		"outコマンドに変換": handler.CreateOutCommands,
		"sheatを保存":   handler.HandleSaveResult,
		"results": func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
			options := i.ApplicationCommandData().Options
			switch options[0].Name {
			case "url":
				handler.HandleResultsUrl(ctx, s, i, repository)
				return
			case "set-url":
				handler.HandleResultsSetURL(ctx, s, i, repository)
				return
			}
		},
		"mogi": func(ctx context.Context, s *discordgo.Session, i *discordgo.InteractionCreate, repository repository.Repository) {
			options := i.ApplicationCommandData().Options
			switch options[0].Name {
			case "list":
				handler.HandleMogiList(ctx, s, i, repository)
				return
			case "set":
				handler.HandleMogiSet(ctx, s, i, repository)
				return
			case "remove":
				handler.HandleMogiRemove(ctx, s, i, repository)
				return
			case "can":
				handler.HandleMogiCan(ctx, s, i, repository)
				return
			case "temp":
				handler.HandleMogiTemp(ctx, s, i, repository)
				return
			case "sub":
				handler.HandleMogiSub(ctx, s, i, repository)
				return
			}
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		ctx := context.Background()
		repository, err := repository.New(ctx)
		if err != nil {
			log.Fatalf("Cannot create repository: %v", err)
			return
		}
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(ctx, s, i, repository)
			}
		case discordgo.InteractionMessageComponent:
			customID := i.MessageComponentData().CustomID
			if customID == string(constant.SQListSelectCustomIDCan) {
				handler.HandleSelect(ctx, s, i, repository)
			} else if customID == string(constant.SQListSelectCustomIDTemp) {
				handler.HandleSelect(ctx, s, i, repository)
			} else if customID == string(constant.SQListSelectCustomIDSub) {
				handler.HandleSelect(ctx, s, i, repository)
			} else if customID == "lounge_name_select" {
				handler.HandleLoungeNameSelect(ctx, s, i, repository)
			} else if strings.HasPrefix(customID, "button_mogi_") {
				handler.HandleMogiButtonClick(ctx, s, i, repository)
			} else if strings.HasPrefix(customID, "button_") {
				handler.HandleClick(ctx, s, i, repository)
			} else if customID == "mogi_remove_select" {
				handler.HandleMogiRemoveSelect(ctx, s, i, repository)
			} else if strings.HasPrefix(customID, "mogi_select_") {
				handler.HandleMogiSelect(ctx, s, i, repository)
			}
		}
	})
}

func init() {
	s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		ctx := context.Background()
		repository, err := repository.New(ctx)
		if err != nil {
			log.Fatalf("Cannot create repository: %v", err)
			return
		}
		handler.HandleLoungeSQInfo(ctx, s, m, repository)
	})
}

func main() {
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	slashCommands := commands.GetCommands()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommands))
	if *RemoveCommands {
		log.Println("Adding commands...")
		for i, v := range slashCommands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			registeredCommands[i] = cmd
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
