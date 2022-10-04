package components

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Meonako/PhoenixManager/commands"
	"github.com/Meonako/PhoenixManager/model-party"

	"github.com/Meonako/go-error"
	"github.com/bwmarrin/discordgo"
)

var ComponentsHandler = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"confirm":        findParty,
	"create-a-party": createParty,
	"leave-a-party":  leaveParty,
}

func findParty(s *discordgo.Session, i *discordgo.InteractionCreate) {
	error.Must(s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}))

	target := findTarget(i)
	party := party.ActiveParties.FindByTarget(target)

	var response *discordgo.WebhookParams
	if party.IsEmpty() {
		response = &discordgo.WebhookParams{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				Embed("No party active at the moment. Do you want to create a party?"),
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Create a party!",
							Style:    discordgo.PrimaryButton,
							CustomID: "create-a-party",
						},
					},
				},
			},
		}
	} else {
		party.Join(getUserID(i))

		for _, player := range party.Players {
			channel, err := s.UserChannelCreate(player)
			if err != nil {
				log.Printf("Failed to create DM channel : %v", err)
				continue
			}

			_, err = s.ChannelMessageSend(
				channel.ID,
				"<@"+getUserID(i)+"> has joined your party!\nWaiting for "+strconv.Itoa(party.MaxPlayer-party.PlayersCount())+" more players!",
			)
			if err != nil {
				log.Printf("Can not send DM message : %v", err)
				continue
			}
		}

		message := "Found a party! Here is a list of member of your party!\n"
		for number, id := range party.Players {
			message += fmt.Sprintf("%v. <@%v>\n", number+1, id)
		}
		message += fmt.Sprintf("Waiting for %v more players to join", party.Target.GetMaxPlayer()-party.PlayersCount())

		response = &discordgo.WebhookParams{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				Embed(message),
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Leave Party",
							Style:    discordgo.DangerButton,
							CustomID: "leave-a-party",
						},
					},
				},
			},
		}
	}

	_, err := s.FollowupMessageCreate(i.Interaction, false, response)
	error.Must(err)
}

func createParty(s *discordgo.Session, i *discordgo.InteractionCreate) {
	target := findTarget(i)
	pt := party.NewParty(target, getUserID(i))

	error.Must(
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags: discordgo.MessageFlagsEphemeral,
				Embeds: []*discordgo.MessageEmbed{
					Embed(
						fmt.Sprintf("Success. Your party have been made!\nWaiting for %v more players to join", pt.Target.GetMaxPlayer()-pt.PlayersCount()),
					),
				},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Disband Party",
								Style:    discordgo.DangerButton,
								CustomID: "leave-a-party",
							},
						},
					},
				},
			},
		}),
	)
}

func leaveParty(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})

	var response *discordgo.WebhookParams
	pt := party.ActiveParties.FindByPlayer(getUserID(i))
	if pt.IsEmpty() {
		response = &discordgo.WebhookParams{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				Embed("You are not in a party!"),
			},
		}
	} else {
		pt.Leave(getUserID(i))
		response = &discordgo.WebhookParams{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				Embed("Done! You have left the team!"),
			},
		}
	}

	_, err := s.FollowupMessageCreate(i.Interaction, false, response)
	error.Must(err)
}

func findTarget(i *discordgo.InteractionCreate) party.TARGET {
	return commands.FindTarget[getUserID(i)]
}

func Embed(msg string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Description: msg,
		Color:       0x0091f9,
	}
}

func getUserID(i *discordgo.InteractionCreate) string {
	if i.User != nil {
		return i.User.ID
	}

	return i.Member.User.ID
}
