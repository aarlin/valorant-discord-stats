package main 

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"math"
	"os"
	"time"
	"flag"
	"os/signal"
	"syscall"
	"strings"
	"github.com/bwmarrin/discordgo"
	"github.com/dariubs/percent"
	"./structures"
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
	// If the message is "ping" reply with "Pong!"
	if strings.Contains(m.Content, "!career") {
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			var json []byte = retrieveBlitzData(nametag)
			var playerStats ValorantStats = parseValorantData(nametag, json)
		
			var careerDamageStats = playerStats.Stats.Overall.Career.DamageStats
		
			var careerHitRateData HitPercentages = calculateHitPercentages(careerDamageStats)
			content := generateDiscordEmbedContent(nametag, playerStats, careerHitRateData, "career")

			embed := structures.NewEmbed().
				SetTitle("Career Statistics").
				AddField(nametag, content).
				SetColor(16582407).MessageEmbed
				
			s.ChannelMessageSendEmbed(m.ChannelID, embed)

			fmt.Printf("%f", careerHitRateData)
		}

	} else if  strings.Contains(m.Content, "!last20") {
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			var json []byte = retrieveBlitzData(nametag)
			var playerStats ValorantStats = parseValorantData(nametag, json)
		
			var last20DamageStats = playerStats.Stats.Overall.Last20.DamageStats
		
			var lastTwentyHitRateData HitPercentages = calculateHitPercentages(last20DamageStats)
			content := generateDiscordEmbedContent(nametag, playerStats, lastTwentyHitRateData, "last20")

			embed := structures.NewEmbed().
				SetTitle("Last 20 Games Statistics").
				AddField(nametag, content).
				SetColor(16582407).MessageEmbed
			
			s.ChannelMessageSendEmbed(m.ChannelID, embed)

			fmt.Printf("%f", lastTwentyHitRateData)
			
		}
	} else if  strings.Contains(m.Content, "!lastgame") {
		if len(strings.Split(m.Content, " ")) > 1 {
			nametag := strings.Split(m.Content, " ")[1]
			var json []byte = retrieveBlitzData(nametag)
			var playerStats ValorantStats = parseValorantData(nametag, json)
		
			matches := retrieveMatches(playerStats.Id)
			matchStats := retrieveMatchStats(playerStats, matches)
		
			embed := structures.NewEmbed().
			SetTitle("Last Game Statistics").
			AddField("Match 1", matchStats).
			// SetImage(mapImage).
			   SetColor(16582407).MessageEmbed
			
			// embed := generateLastGameStatsEmbed(nametag)
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	} else if m.Content == "!commands" {
		s.ChannelMessageSend(m.ChannelID, "Commands are\n!career <nametag>\n!last20 <nametag>\n!lastgame <nametag\n")
	}
}


type HitPercentages struct {
	HeadShotPercentage float64
	BodyShotPercentage float64
	LegShotPercentage float64
}

type DamageStats struct {
	Bodyshots int
	Headshots int
	Legshots  int
	Damage    int
}
type MatchStatistics struct {
	DamageStats DamageStats
	Matches int
}
type Overall struct {
	Career MatchStatistics
	Last20 MatchStatistics
}
type Stats struct {
	Overall Overall
}
type ValorantStats struct {
	Nametag string
	Id string
	Stats Stats	
}

type Damage struct {
	Receiver  string `json:"receiver"`
	Damage    int    `json:"damage"`
	Legshots  int    `json:"legshots"`
	Bodyshots int    `json:"bodyshots"`
	Headshots int    `json:"headshots"`
}

type PlayerStat struct {
	Score  int `json:"score"`
	Damage []Damage `json:"damage"`
	WasAfk  bool `json:"wasAfk"`
	Subject      string `json:"subject"`
	WasPenalized bool   `json:"wasPenalized"`
}

type RoundResult struct { 
	RoundNum    int    `json:"roundNum"`
	PlantSite   string `json:"plantSite"`
	BombPlanter string `json:"bombPlanter,omitempty"`
	PlayerStats []PlayerStat `json:"playerStats"`
	RoundResult  string `json:"roundResult"`
	WinningTeam  string `json:"winningTeam"`
}

type MatchHistory struct {
	ID     		 string 	   `json:"id"`
	Map    	     string 	   `json:"map"`
	Mode   		 string 	   `json:"mode"`
	Ranked 		 bool   	   `json:"ranked"`
	RoundResults []RoundResult `json:"roundResults"`
	StartedAt 	 time.Time 	   `json:"startedAt"`
	Length    	 int       	   `json:"length"`
	Queue        string        `json:"queue"`
	Season    	 string        `json:"season"`
	Version   	 string        `json:"version"`
}

type MatchHistoryOffset struct {
	Count int
	Data []MatchHistory
	Limit int
	Offset string
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

func retrieveBlitzData(nametag string) []byte {
	blitzEndpoint := fmt.Sprintf("https://valorant.iesdev.com/player/%s", nametag)
	resp, err := http.Get(blitzEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func parseValorantData(nametag string, blitzJson []byte) ValorantStats {
	var valorantStats ValorantStats 
	err := json.Unmarshal(blitzJson, &valorantStats)
	if err != nil {
		postError(nametag)
		log.Fatalln(err)
	}
	return valorantStats
}

func retrieveMatches(playerID string) []string {
	matchHistoryEndpoint := fmt.Sprintf("https://valorant.iesdev.com/matches/%s?offset=0&queue=", playerID)
	resp, err := http.Get(matchHistoryEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var matchHistoryOffset MatchHistoryOffset 
	errUnmarshal := json.Unmarshal(body, &matchHistoryOffset) 
	if errUnmarshal != nil {
		postError(playerID)
		log.Fatalln(errUnmarshal)
	}

	var matches []string = make([]string, 0, 10)

	for _, match := range matchHistoryOffset.Data {
		matches = append(matches, match.ID)
	}

	// fmt.Println(matches)
	return matches

}

func retrieveMatchStats(player ValorantStats, matches []string) string {
	var matchDamageStatistics = make(map[string]*DamageStats)
	
	matchEndpoint := fmt.Sprintf("https://valorant.iesdev.com/match/%s", matches[0])
	resp, err := http.Get(matchEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var matchHistory MatchHistory 
	errUnmarshal := json.Unmarshal(body, &matchHistory) 
	if errUnmarshal != nil {
		postError(player.Nametag)
		log.Fatalln(errUnmarshal)
	}

	// make call to each match id using get request 
	// then do this nested for loop to grab all damage
	matchDamageStatistics[matchHistory.ID] = &DamageStats{0, 0, 0, 0}
	for _, roundResult := range matchHistory.RoundResults {
		for _, playerStat := range roundResult.PlayerStats {
			if playerStat.Subject == player.Id {
				for _, damage := range playerStat.Damage { // array of damages done to different enemies in the round
					matchDamageStatistics[matchHistory.ID].Damage += damage.Damage
					matchDamageStatistics[matchHistory.ID].Headshots += damage.Headshots
					matchDamageStatistics[matchHistory.ID].Bodyshots += damage.Bodyshots
					matchDamageStatistics[matchHistory.ID].Legshots += damage.Legshots
				}
			}
		}
	}
	
	var matchPercentages = calculateHitPercentages(*matchDamageStatistics[matchHistory.ID])
	fmt.Println(matchPercentages)

	var matchStats = fmt.Sprintf("Nametag: %s\nMatch ID: %s\nMap: %s\nHeadshots: %d (%.2f%%)\nBodyshots: %d (%.2f%%)\nLegshots: %d(%.2f%%)\nDamage: %d\n", 
		player.Nametag,
		matchHistory.ID, 
		matchHistory.Map,
		matchDamageStatistics[matchHistory.ID].Headshots, matchPercentages.HeadShotPercentage,
		matchDamageStatistics[matchHistory.ID].Bodyshots, matchPercentages.BodyShotPercentage,
		matchDamageStatistics[matchHistory.ID].Legshots, matchPercentages.LegShotPercentage,
		matchDamageStatistics[matchHistory.ID].Damage)
	fmt.Println(matchStats)

	// var mapImage = fmt.Sprintf("https://blitz-cdn.blitz.gg/blitz/val/maps/map-art-%s.jpg", matchHistory.Map)

	return matchStats
}


func calculateHitPercentages(damageStats DamageStats) HitPercentages {
	headShots := damageStats.Headshots
	bodyShots := damageStats.Bodyshots
	legShots := damageStats.Legshots

	var hitPercentages HitPercentages
	totalShots := headShots + bodyShots + legShots
	hitPercentages.HeadShotPercentage = roundPercentage(percent.PercentOf(headShots, totalShots))
	hitPercentages.BodyShotPercentage = roundPercentage(percent.PercentOf(bodyShots, totalShots))
	hitPercentages.LegShotPercentage  = roundPercentage(percent.PercentOf(legShots, totalShots))
	
	return hitPercentages
}

func roundPercentage(percentage float64) float64 {
	return math.Round(percentage * 100) / 100
}

func postError(nametag string) {
	content := fmt.Sprintf("Could not retrieve data for %s. Check if you linked blitz.gg with your account.", nametag)
	discordWebhook := "https://discordapp.com/api/webhooks/723323733728821369/amDzaBkpO80fWYPJbRejem39CSa00zRdFcF4SO5tYMtprP3V8vsT6autU3nG3ik9TOuc"
	discordMessage := map[string]interface{} {
		"content": content,
	}
	
	bytesRepresentation, err := json.Marshal(discordMessage)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(discordWebhook, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result)
	log.Println(result["data"])
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
