package main 

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"
	"flag"
	"os/signal"
	"syscall"
	"errors"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/dariubs/percent"
	"github.com/aarlin/valorant-discord-stats/structures"
	"github.com/aarlin/valorant-discord-stats/handler"
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
		// TODO: Add K/D ratio, dmg/round, combat score, econ score, current rank
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			blitzData, err := handler.RetrieveBlitzData(nametag)
			if (err != nil) {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if (err != nil) {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					var careerDamageStats = playerStats.Stats.Overall.Career.DamageStats
				
					var careerHitRateData HitPercentages = calculateHitPercentages(careerDamageStats)
					content := generateDiscordEmbedContent(nametag, playerStats, careerHitRateData, "career")
		
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
			blitzData, err := handler.RetrieveBlitzData(nametag)
			if (err != nil) {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if (err != nil) {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					var last20DamageStats = playerStats.Stats.Overall.Last20.DamageStats
				
					var lastTwentyHitRateData HitPercentages = calculateHitPercentages(last20DamageStats)
					content := generateDiscordEmbedContent(nametag, playerStats, lastTwentyHitRateData, "last20")
		
					embed := structures.NewEmbed().
						SetTitle("Last 20 Games Statistics").
						AddField(nametag, content).
						SetColor(16582407).MessageEmbed
					
					s.ChannelMessageSendEmbed(m.ChannelID, embed)
		
					fmt.Printf("%f\n", lastTwentyHitRateData)
				}
			}

		}
	} else if strings.Contains(m.Content, "!lastgame") {
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			blitzData, err := handler.RetrieveBlitzData(nametag)
			if (err != nil) {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if (err != nil) {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					matches, err := handler.RetrieveMatches(playerStats.Id)
					if (err != nil) {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						matchStats, err := handler.RetrieveMatchStats(playerStats, matches)
						
						if (err != nil) {
							s.ChannelMessageSend(m.ChannelID, err.Error())
						} else {
							embed := structures.NewEmbed().
							SetTitle("Last Game Statistics").
							AddField("Match 1", matchStats).
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
	
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func parseValorantData(nametag string, blitzJson []byte) (ValorantStats, error) {
	var valorantStats ValorantStats 
	err := json.Unmarshal(blitzJson, &valorantStats)
	if err != nil {
		retrieveDataErr := fmt.Sprintf("Could not retrieve data for %s. Check if you linked blitz.gg with your account.", nametag)
		return valorantStats, errors.New(retrieveDataErr)
	}
	return valorantStats, nil
}

func createCompetitiveTier(competitiveTier int) string {
	// TODO: Create mapping for competitive tiers
	return ""
}

func generateDiscordEmbedContent(nametag string, stats ValorantStats, hitRate HitPercentages, matchStatisticType string) string {
	headShots := fmt.Sprintf(":no_mouth: Head shot percentage: %.2f%%\n", hitRate.HeadShotPercentage)
	bodyShots := fmt.Sprintf(":shirt: Body shot percentage: %.2f%%\n", hitRate.BodyShotPercentage)
	legShots := fmt.Sprintf(":foot: Leg shot percentage: %.2f%%\n", hitRate.LegShotPercentage)
	matchesPlayed := stats.Stats.Overall.Career.Matches
	content := ""
	switch matchStatisticType {
		case "career": 	content = fmt.Sprintf("Career Stats for %s:\nTotal number of matches: %d\n%s%s%s\n", nametag, matchesPlayed, headShots, bodyShots, legShots)
		case "last20": content = fmt.Sprintf("Last 20 Games Stats for %s:\n%s%s%s\n", nametag, headShots, bodyShots, legShots)
		default: content = fmt.Sprintf("Something went wrong choosing the statistic type when posting to Discord")
	}

	return content
}
