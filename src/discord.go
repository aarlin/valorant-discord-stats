package main 

import (
	"log"
	"net/http"
	"encoding/json"
	"bytes"
	"fmt"
)

func main() {
	retrieveBlitzData("fompei-na1")
}

func retrieveBlitzData(nametag string) string {
	blitzEndpoint := fmt.Sprintf("https://valorant.iesdev.com/player/%s", nametag)
	resp, err := http.Get(blitzEndpoint)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	type DamageStats struct {
		bodyShots int
		headShots int
		legShots int
		damage int
	}
	type Career struct {
		damageStats DamageStats
	}
	type Overall struct {
		career Career
	}
	type Stats struct {
		overall Overall
	}
	type ValorantStats struct {
		stats Stats
	}

	var result map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&result)
	log.Println(result)
	return ""
}

func postStatsToDiscord() {
	discordWebhook := "https://discordapp.com/api/webhooks/723323733728821369/amDzaBkpO80fWYPJbRejem39CSa00zRdFcF4SO5tYMtprP3V8vsT6autU3nG3ik9TOuc"
	discordMessage := map[string]interface{} {
		"content": "Valorant discord post",
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