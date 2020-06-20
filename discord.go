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
type Career struct {
	DamageStats DamageStats
	Matches int
}
type Overall struct {
	Career Career
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

	nametag := "fompei-na1"
	json := retrieveBlitzData(nametag)
	playerStats := parseValorantData(nametag, json)
	hitRateData := calculateHitPercentages(playerStats)
	postStatsToDiscord(nametag, playerStats, hitRateData)
	fmt.Println(hitRateData)
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

func calculateHitPercentages(valorantStats ValorantStats) HitPercentages {
	headShots := valorantStats.Stats.Overall.Career.DamageStats.HeadShots
	bodyShots := valorantStats.Stats.Overall.Career.DamageStats.BodyShots
	legShots  := valorantStats.Stats.Overall.Career.DamageStats.LegShots

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

func postStatsToDiscord(nametag string, stats ValorantStats, hitRate HitPercentages) {
	headShots := fmt.Sprintf(":no_mouth: Head shot percentage: %.2f%%\n", hitRate.HeadShotPercentage)
	bodyShots := fmt.Sprintf(":shirt: Body shot percentage: %.2f%%\n", hitRate.BodyShotPercentage)
	legShots := fmt.Sprintf(":foot: Leg shot percentage: %.2f%%\n", hitRate.LegShotPercentage)
	matchesPlayed := stats.Stats.Overall.Career.Matches
	content := fmt.Sprintf("Career Stats for %s:\n%s%s%s\n Total number of matches: %d", nametag, headShots, bodyShots, legShots, matchesPlayed)

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
