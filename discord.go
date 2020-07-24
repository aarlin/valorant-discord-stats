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
			blitzData, err := retrieveBlitzData(nametag)
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
			blitzData, err := retrieveBlitzData(nametag)
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
			blitzData, err := retrieveBlitzData(nametag)
			if (err != nil) {
				s.ChannelMessageSend(m.ChannelID, err.Error())
			} else {
				playerStats, err := parseValorantData(nametag, blitzData)
				if (err != nil) {
					s.ChannelMessageSend(m.ChannelID, err.Error())
				} else {
					matches, err := retrieveMatches(playerStats.Id)
					if (err != nil) {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						matchStats, err := retrieveMatchStats(playerStats, matches)
						
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

type Ranks struct {
	Competitive struct {
		Tier int 
	} `json:"competitive"`
}

type ValorantStats struct {
	Nametag string
	Id string
	Stats Stats	
	Ranks Ranks
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

type Team struct {
	Won          bool   `json:"won"`
	TeamID       string `json:"teamId"`
	RoundsWon    int    `json:"roundsWon"`
	RoundsPlayed int    `json:"roundsPlayed"`
}

type RoundResult struct { 
	RoundNum    int    `json:"roundNum"`
	PlantSite   string `json:"plantSite"`
	BombPlanter string `json:"bombPlanter,omitempty"`
	PlayerStats []PlayerStat `json:"playerStats"`
	RoundResult  string `json:"roundResult"`
	WinningTeam  string `json:"winningTeam"`
}

type PlayerStats struct {
	Kills        int `json:"kills"`
	Score        int `json:"score"`
	Deaths       int `json:"deaths"`
	Assists      int `json:"assists"`
	AbilityCasts struct {
		GrenadeCasts  int `json:"grenadeCasts"`
		Ability1Casts int `json:"ability1Casts"`
		Ability2Casts int `json:"ability2Casts"`
		UltimateCasts int `json:"ultimateCasts"`
	} `json:"abilityCasts"`
	RoundsPlayed   int `json:"roundsPlayed"`
	PlaytimeMillis int `json:"playtimeMillis"`
}

type Player struct {
	Stats  					PlayerStats	`json:"stats"`
	TeamID      			string 		`json:"teamId"`			// blue or red
	PartyID    	 			string 		`json:"partyId"`
	Subject     			string 		`json:"subject"`			// player id
	CharacterID 			string 		`json:"characterId"`		// agent
	CompetitiveTier         int 		`json:"competitiveTier"`
	SessionPlaytimeMinutes  int 		`json:"sessionPlaytimeMinutes"`
	RoundDamage []struct {
		Round    int    `json:"round"`
		Damage   int    `json:"damage"`
		Receiver string `json:"receiver"`
	} `json:"roundDamage"`
}

type MatchHistory struct {
	ID     		 string 	   `json:"id"`
	Map    	     string 	   `json:"map"`
	Mode   		 string 	   `json:"mode"`
	Ranked 		 bool   	   `json:"ranked"`
	Teams  		 []Team 	   `json:"teams"`
	RoundResults []RoundResult `json:"roundResults"`
	StartedAt 	 time.Time 	   `json:"startedAt"`
	Players 	 []Player      `json:"players"`
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

func retrieveBlitzData(nametag string) ([]byte, error) {
	blitzEndpoint := fmt.Sprintf("https://valorant.iesdev.com/player/%s", nametag)
	resp, err := http.Get(blitzEndpoint)
	if err != nil {
		return nil, errors.New("Unable to get data from player endpoint")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Unable to read response body")
	}

	return body, nil
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

func retrieveMatches(playerID string) ([]string, error) {
	matchHistoryEndpoint := fmt.Sprintf("https://valorant.iesdev.com/matches/%s?offset=0&queue=", playerID)
	resp, err := http.Get(matchHistoryEndpoint)
	if err != nil {
		return nil, errors.New("Unable to get data from matches endpoint")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Unable to read response body")
	}

	var matchHistoryOffset MatchHistoryOffset 
	errUnmarshal := json.Unmarshal(body, &matchHistoryOffset) 
	if errUnmarshal != nil {
		return nil, errors.New("Error trying to parse matches in history")
	}

	var matches []string = make([]string, 0, 10)

	for _, match := range matchHistoryOffset.Data {
		matches = append(matches, match.ID)
	}

	// fmt.Println(matches)
	return matches, nil

}

func retrieveMatchStats(player ValorantStats, matches []string) (string, error) {
	var matchDamageStatistics = make(map[string]*DamageStats)
	
	matchEndpoint := fmt.Sprintf("https://valorant.iesdev.com/match/%s", matches[0])
	resp, err := http.Get(matchEndpoint)
	if err != nil {
		return "", errors.New("Unable to get data from match endpoint")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Unable to read response body")
	}

	var matchHistory MatchHistory 
	errUnmarshal := json.Unmarshal(body, &matchHistory) 
	if errUnmarshal != nil {
		retrieveDataErr := fmt.Sprintf("Could not retrieve data for %s. Check if you linked blitz.gg with your account.", player.Nametag)
		return "", errors.New(retrieveDataErr)
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

	var kills int
	var deaths int
	var assists int
	var competitiveTier int
	var score int
	var roundsPlayed int
	var team string 
	// TODO: parse round damage for dmg/round
	// var roundDamage struct
	gameRoundResults := "0 - 0"
	fmt.Println(team)

	for _, matchParticipant := range matchHistory.Players {
		if matchParticipant.Subject == player.Id {
			kills = matchParticipant.Stats.Kills
			deaths = matchParticipant.Stats.Deaths
			assists = matchParticipant.Stats.Assists
			competitiveTier = matchParticipant.CompetitiveTier
			score = matchParticipant.Stats.Score
			roundsPlayed = matchParticipant.Stats.RoundsPlayed
			team = matchParticipant.TeamID
			fmt.Println(team)

			var roundDamageMap = make(map[int]int)
			for _, round := range matchParticipant.RoundDamage {
				fmt.Printf("round")
				fmt.Println(round)
				if val, ok := roundDamageMap[round.Round]; ok {
					fmt.Printf("val")
					fmt.Println(val)
					// mapping is correct BUT missing rounds where player does no damage
					roundDamageMap[round.Round] = val + round.Damage
				} else {
					roundDamageMap[round.Round] = round.Damage
				}
			}
			fmt.Println(roundDamageMap)
		}
	}

	var blueTeamRoundWins int
	var redTeamRoundWins int

	for _, matchParticipantTeam := range matchHistory.Teams {
		switch matchParticipantTeam.TeamID {
			case "Red": redTeamRoundWins = matchParticipantTeam.RoundsWon
			case "Blue": blueTeamRoundWins = matchParticipantTeam.RoundsWon
			default: fmt.Sprintf("Something went wrong parsing team ID from API endpoint")
		}
	}
	fmt.Println(team)

	switch team {
		case "Red":	gameRoundResults = fmt.Sprintf("%d - %d", redTeamRoundWins, blueTeamRoundWins)
		case "Blue": gameRoundResults = fmt.Sprintf("%d - %d", blueTeamRoundWins, redTeamRoundWins)
		default: fmt.Sprintf("Something went wrong while trying to create round results")
	}
	
	var matchPercentages = calculateHitPercentages(*matchDamageStatistics[matchHistory.ID])
	fmt.Println(matchPercentages)

	var matchStats = fmt.Sprintf("Nametag: %s\n" + 
		"Competitive Tier: %s\n" + 
		"Game Results: %s\n" + 
		"Match ID: %s\n" +
		"Map: %s\n" + 
		"Headshots: %d (%.2f%%)\n" +
		"Bodyshots: %d (%.2f%%)\n" + 
		"Legshots: %d (%.2f%%)\n" + 
		"Damage: %d\n" + 
		"Combat Score: %d\n" + 
		"K\\/D\\/A: %d\\/%d\\/%d\n",
		player.Nametag,
		competitiveTier,
		gameRoundResults,
		matchHistory.ID, 
		matchHistory.Map,
		matchDamageStatistics[matchHistory.ID].Headshots, matchPercentages.HeadShotPercentage,
		matchDamageStatistics[matchHistory.ID].Bodyshots, matchPercentages.BodyShotPercentage,
		matchDamageStatistics[matchHistory.ID].Legshots, matchPercentages.LegShotPercentage,
		matchDamageStatistics[matchHistory.ID].Damage,
		(score / roundsPlayed),
		kills, deaths, assists)
	fmt.Println(matchStats)

	// TODO: Add map image 
	// TODO: Add competitive tier image
	// var mapImage = fmt.Sprintf("https://blitz-cdn.blitz.gg/blitz/val/maps/map-art-%s.jpg", matchHistory.Map)

	return matchStats, nil
}

func createCompetitiveTier(competitiveTier int) string {
	// TODO: Create mapping for competitive tiers
	return ""
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
