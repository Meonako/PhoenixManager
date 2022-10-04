package commands

import (
	"strings"

	party "github.com/Meonako/PhoenixManager/model-party"
	"github.com/bwmarrin/discordgo"
)

var CommandsList []*discordgo.ApplicationCommand = []*discordgo.ApplicationCommand{
	{
		Name:        "join",
		Options:     generateJoinCommandOptions(),
		Description: "Find a party for you :)",
	},
	{
		Name:        "leave",
		Description: "Leave party if you're in one",
	},
	{
		Name:        "clear",
		Description: "Clear 100 messages in channel",
	},
	{
		Name: "crit",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "critical",
				Description: "Number of your CRIT",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "enemy-level",
				Description: "Critical chance is base on Enemy Level",
				Required:    true,
			},
		},
		Description: "Calculate your critical chance",
	},
}

// Generate commands
func generateJoinCommandOptions() (sub []*discordgo.ApplicationCommandOption) {
	for i := 0; i < party.Len; i++ {
		mode := party.TARGET(i)
		cmd := &discordgo.ApplicationCommandOption{
			Type: discordgo.ApplicationCommandOptionSubCommand,
			Name: strings.ReplaceAll(
				strings.ToLower(mode.String()), " ", "-",
			),
			Description: mode.String() + " Party",
		}
		sub = append(sub, cmd)
	}
	return
}
