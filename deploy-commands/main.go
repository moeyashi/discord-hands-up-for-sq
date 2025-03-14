package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/moeyashi/discord-hands-up-for-sq/commands"
)

func main() {
	EnvBotToken := os.Getenv("BOT_TOKEN")
	s, err := discordgo.New("Bot " + EnvBotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	defer s.Close()
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Removing commands...")
	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}
	// for _, v := range registeredCommands {
	// 	err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
	// 	if err != nil {
	// 		log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
	// 	}
	// }

	log.Println("Adding commands...")
	for _, v := range commands.GetCommands() {
		existsCommand := false
		for _, prev := range registeredCommands {
			if prev.Name == v.Name {
				existsCommand = true
				log.Println("Editing command:", v.Name)
				_, err := s.ApplicationCommandEdit(s.State.User.ID, "", prev.ID, v)
				if err != nil {
					log.Printf("Cannot edit '%v' command: %v", v.Name, err)
				}
				break
			}
		}

		if existsCommand {
			continue
		}

		log.Println("Creating command:", v.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}
