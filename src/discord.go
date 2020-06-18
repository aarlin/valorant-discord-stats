package main 

import (
	"log"
	"net/http"
	"encoding/json"
	"bytes"
)

func main() {
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