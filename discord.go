package main 

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
	"math"
	"github.com/dariubs/percent"
	"os"
	"time"
)

type HitPercentages struct {
	HeadShotPercentage float64
	BodyShotPercentage float64
	LegShotPercentage float64
}

type DamageStats struct {
	BodyShots int
	HeadShots int
	LegShots  int
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
	Damage    int    `json:"damage"`
	Legshots  int    `json:"legshots"`
	Receiver  string `json:"receiver"`
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
	nametag := os.Args[1]
	var json []byte = retrieveBlitzData(nametag)
	var playerStats ValorantStats = parseValorantData(nametag, json)
	var careerHitRateData HitPercentages = calculateHitPercentages(playerStats, "career")
	var lastTwentyHitRateData HitPercentages = calculateHitPercentages(playerStats, "last20")
	// postStatsToDiscord(nametag, playerStats, careerHitRateData, "career")
	// postStatsToDiscord(nametag, playerStats, lastTwentyHitRateData, "last20")
	postToSpreadsheet(playerStats.Id)
	fmt.Printf("%f%f", careerHitRateData, lastTwentyHitRateData)
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

func postToSpreadsheet(playerID string) {
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

	for index, match := range matchHistoryOffset.Data {
		fmt.Printf("%d: %s\n", index, match.ID)
	}

}

func calculateHitPercentages(valorantStats ValorantStats, matchStatisticType string) HitPercentages {
	headShots := 0
	bodyShots := 0
	legShots := 0

	switch matchStatisticType {
		case "career": 		
			headShots = valorantStats.Stats.Overall.Career.DamageStats.HeadShots
			bodyShots = valorantStats.Stats.Overall.Career.DamageStats.BodyShots
			legShots  = valorantStats.Stats.Overall.Career.DamageStats.LegShots
		case "last20":
			headShots = valorantStats.Stats.Overall.Last20.DamageStats.HeadShots
			bodyShots = valorantStats.Stats.Overall.Last20.DamageStats.BodyShots
			legShots  = valorantStats.Stats.Overall.Last20.DamageStats.LegShots

		default: fmt.Println("Incorrect matchStatisticType, please choose between career and last20")
	}

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

func postStatsToDiscord(nametag string, stats ValorantStats, hitRate HitPercentages, matchStatisticType string) {
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

	// generated from https://mholt.github.io/json-to-go/

	type Field struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
	type Embed struct {
		Title  string  `json:"title"`
		Color  int     `json:"color"`
		Fields []Field `json:"fields"`
	}
	type DiscordMessage struct {
		Embeds []Embed `json:"embeds"`
	}	


	discordWebhook := "https://discordapp.com/api/webhooks/723323733728821369/amDzaBkpO80fWYPJbRejem39CSa00zRdFcF4SO5tYMtprP3V8vsT6autU3nG3ik9TOuc"
	discordMessage := DiscordMessage{
		Embeds: []Embed{
			Embed{
				Title: "Valorant Statistics",
				Color: 16582407,
				Fields: []Field{
					Field{
						Name: "Career Statistics",
						Value: content,
					},
				},
			},
		},
	}

	fmt.Println(discordMessage)

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
