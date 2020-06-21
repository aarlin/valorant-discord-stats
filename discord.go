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


func main() {
	nametag := os.Args[1]
	json := retrieveBlitzData(nametag)
	playerStats := parseValorantData(nametag, json)
	careerHitRateData := calculateHitPercentages(playerStats, "career")
	lastTwentyHitRateData := calculateHitPercentages(playerStats, "last20")
	postStatsToDiscord(nametag, playerStats, careerHitRateData, "career")
	postStatsToDiscord(nametag, playerStats, lastTwentyHitRateData, "last20")
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
