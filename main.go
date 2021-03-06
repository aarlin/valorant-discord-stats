package main

import (
	"fmt"
	"os"
	"log"
	"os/signal"
	"strings"
	"syscall"
	"github.com/aarlin/valorant-discord-stats/src/definition"
	"github.com/aarlin/valorant-discord-stats/src/api"
	"github.com/aarlin/valorant-discord-stats/src/calculation"
	"github.com/aarlin/valorant-discord-stats/src/formatter"
	"github.com/aarlin/valorant-discord-stats/src/structures"
	"github.com/aarlin/valorant-discord-stats/src/config"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

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
				playerStats, err := api.ParseValorantData(nametag, blitzData)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					var careerDamageStats = playerStats.Stats.Overall.Career.DamageStats

					var careerHitRateData definition.HitPercentages = calculation.CalculateHitPercentages(careerDamageStats)
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
				playerStats, err := api.ParseValorantData(nametag, blitzData)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					var last20DamageStats = playerStats.Stats.Overall.Last20.DamageStats

					var lastTwentyHitRateData definition.HitPercentages = calculation.CalculateHitPercentages(last20DamageStats)
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
				playerStats, err := api.ParseValorantData(nametag, blitzData)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					matches, err := api.RetrieveMatches(playerStats.Id)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						matchHistory, err := api.RetrieveMatchHistory(playerStats, matches)
						matchSummary, matchSummaryText := formatter.GenerateMatchSummary(playerStats, matchHistory)


						if err != nil {
							s.ChannelMessageSend(m.ChannelID, err.Error())
						} else {
							embed := structures.NewEmbed().
								SetTitle("Last Game Statistics").
								SetURL(formatter.GenerateMatchLink(nametag, matchSummary.ID)).
								SetThumbnail(formatter.GenerateMapImageLink(matchSummary.Map)).
								SetFooter(formatter.GenerateCompetitiveTierFooter(matchSummary.CompetitiveTier)).
								SetColor(config.EMBED_COLOR).
								AddField(matchSummary.ID, matchSummaryText).MessageEmbed
						
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
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
  
	discordToken := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + discordToken)
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