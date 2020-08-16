package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"github.com/aarlin/valorant-discord-stats/definition"
	"github.com/aarlin/valorant-discord-stats/api"
	"github.com/aarlin/valorant-discord-stats/calculation"
	"github.com/aarlin/valorant-discord-stats/formatter"
	"github.com/aarlin/valorant-discord-stats/structures"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Contains(m.Content, "!career") {
		// TODO: econ score, current rank
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			blitzData, err := api.RetrieveBlitzData(nametag)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					var careerDamageStats = playerStats.Stats.Overall.Career.DamageStats

					var careerHitRateData definition.HitPercentages = helper.CalculateHitPercentages(careerDamageStats)
					content := formatter.GenerateDiscordEmbedContent(nametag, playerStats, careerHitRateData, "career")

					embed := structures.NewEmbed().
						SetTitle("Career Statistics").
						AddField(nametag, content).
						SetColor(16582407).MessageEmbed

					s.ChannelMessageSendEmbed(m.ChannelID, embed)

					fmt.Printf("%f\n", careerHitRateData)
				}
			}
		}

	} else if strings.Contains(m.Content, "!last20") {
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			blitzData, err := api.RetrieveBlitzData(nametag)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					var last20DamageStats = playerStats.Stats.Overall.Last20.DamageStats

					var lastTwentyHitRateData definition.HitPercentages = helper.CalculateHitPercentages(last20DamageStats)
					content := formatter.GenerateDiscordEmbedContent(nametag, playerStats, lastTwentyHitRateData, "last20")

					embed := structures.NewEmbed().
						SetTitle("Last 20 Games Statistics").
						AddField(nametag, content).
						SetColor(16582407).MessageEmbed

					s.ChannelMessageSendEmbed(m.ChannelID, embed)

					fmt.Printf("%f\n", lastTwentyHitRateData)
				}
			}

		}
	} else if strings.Contains(m.Content, "!lastgame") || strings.Contains(m.Content, "!lg") {
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			blitzData, err := api.RetrieveBlitzData(nametag)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					matches, err := api.RetrieveMatches(playerStats.Id)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						matchStats, err := api.RetrieveMatchStats(playerStats, matches)
						matchSummary, err := api.GenerateMatchSummary(playerStats, matchStats)

						if err != nil {
							s.ChannelMessageSend(m.ChannelID, err.Error())
						} else {
							embed := structures.NewEmbed().
								AddField("Last Game Statistics", matchSummary).
								// SetImage(mapImage).
								SetColor(16582407).MessageEmbed

							// embed := generateLastGameStatsEmbed(nametag)
							s.ChannelMessageSendEmbed(m.ChannelID, embed)
						}
					}
				}
			}
		}
	} else if m.Content == "!commands" {
		s.ChannelMessageSend(m.ChannelID, "Commands are\n!career <nametag>\n!last20 <nametag>\n!lastgame <nametag\n")
	}
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func parseValorantData(nametag string, blitzJson []byte) (definition.ValorantStats, error) {
	var valorantStats definition.ValorantStats
	err := json.Unmarshal(blitzJson, &valorantStats)
	if err != nil {
		retrieveDataErr := fmt.Sprintf("Could not retrieve data for %s. Check if you linked blitz.gg with your account.", nametag)
		return valorantStats, errors.New(retrieveDataErr)
	}
	return valorantStats, nil
}