package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/Meonako/PhoenixManager/commands"
	"github.com/Meonako/PhoenixManager/components"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token string
	AppID string

	session *discordgo.Session
)

func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	Token = os.Getenv("BotToken")
	AppID = os.Getenv("AppID")

	if Token == "" {
		log.Fatal("No bot token provided")
	}

	if AppID == "" {
		log.Fatal("No app ID provided")
	}

	session, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
}

func main() {
	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			var commander *discordgo.User // IGNORE THIS VARIABLE NAME. ITS FOR THE MEMES
			if i.User != nil {
				commander = i.User
			} else {
				commander = i.Member.User
			}
			log.Printf("%v#%v issued : %v", commander.Username, commander.Discriminator, i.ApplicationCommandData().Name)

			if h, ok := commands.CommandsHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := components.ComponentsHandler[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	for _, command := range commands.CommandsList {
		_, err := session.ApplicationCommandCreate(AppID, "", command)
		if err != nil {
			log.Panicf("Cannot create slash command: %v", err)
		}
		log.Printf("Registered command: %v\n", command.Name)
	}

	err := session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Ready! Press CTRL + C to exit properly!")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Cleaning up...")
	registeredCommands, err := session.ApplicationCommands(AppID, "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v\n", err)
	}

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(AppID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}

		log.Printf("Deleted '%v' command\n", v.Name)
	}
}
