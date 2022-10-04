package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Meonako/PhoenixManager/model-party"

	"github.com/Meonako/go-error"
	"github.com/bwmarrin/discordgo"
)

var (
	CommandsHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"join":  join,
		"leave": leave,
		"clear": clear,
		"crit":  crit,
	}

	FindTarget = map[string]party.TARGET{}
)

func join(s *discordgo.Session, i *discordgo.InteractionCreate) {
	saveTarget(i)
	error.Must(s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Do you want to TRY to join **" +
				strings.ToUpper(party.TARGET(
					party.GetEnum(
						strings.ToUpper(i.ApplicationCommandData().Options[0].Name),
					)).String(),
				) + "** party?",
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Yes",
							Style:    discordgo.SuccessButton,
							CustomID: "confirm",
						},
						discordgo.Button{
							Label:    "No",
							Style:    discordgo.DangerButton,
							CustomID: "fc_no",
						},
					},
				},
			},
		},
	}))
}

func leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				createEmbed("Are you sure you want to leave your current party?"),
			},
			Flags: discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Confirm",
							Style:    discordgo.SuccessButton,
							CustomID: "leave-a-party",
						},
					},
				},
			},
		},
	})
}

func clear(s *discordgo.Session, i *discordgo.InteractionCreate) {
	error.Must(
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}),
	)

	message, err := s.ChannelMessages(i.ChannelID, 100, "", "", "")
	error.Must(err)

	for _, msg := range message {
		s.ChannelMessageDelete(i.ChannelID, msg.ID)
	}

	_, err = s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			createEmbed("Deleted " + strconv.Itoa(len(message)) + " Messages!"),
		},
		Flags: discordgo.MessageFlagsEphemeral,
	})
	error.Must(err)
}

func crit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	error.Must(
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
			},
		}),
	)

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	critical := float64(optionMap["critical"].IntValue())
	enemyLevel := float64(optionMap["enemy-level"].IntValue())
	critChance := critical / (enemyLevel*2.666 - 27)

	_, err := s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: "Your critical rate against LV. " + strconv.Itoa(int(enemyLevel)) + " enemies is " + fmt.Sprintf("%f", critChance) + "%",
		Flags:   discordgo.MessageFlagsEphemeral,
	})
	error.Must(err)
}

func createEmbed(msg string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: msg,
		Color:       0x0091f9,
	}
}

func saveTarget(i *discordgo.InteractionCreate) {
	target := party.GetEnum(strings.ToUpper(i.ApplicationCommandData().Options[0].Name))
	FindTarget[getUserID(i)] = party.TARGET(target)
}

func getUserID(i *discordgo.InteractionCreate) string {
	if i.User != nil {
		return i.User.ID
	}

	return i.Member.User.ID
}
