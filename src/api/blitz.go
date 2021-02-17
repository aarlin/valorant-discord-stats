package api

import (
	"fmt"
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
	"github.com/aarlin/valorant-discord-stats/src/definition"
)

func RetrieveBlitzData(nametag string) ([]byte, error) {
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

func RetrieveMatches(playerID string) ([]string, error) {
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

	var matchHistoryOffset definition.MatchHistoryOffset 
	errUnmarshal := json.Unmarshal(body, &matchHistoryOffset) 
	if errUnmarshal != nil {
		return nil, errors.New("Error trying to parse matches in history")
	}

	var matches []string = make([]string, 0, 10)

	for _, match := range matchHistoryOffset.Data {
		matches = append(matches, match.ID)
	}

	return matches, nil

}

func RetrieveMatchHistory(player definition.ValorantStats, matches []string) (definition.MatchHistory, error) {
	var matchHistory definition.MatchHistory 

	matchEndpoint := fmt.Sprintf("https://valorant.iesdev.com/match/%s", matches[0])
	resp, err := http.Get(matchEndpoint)
	if err != nil {
		return matchHistory, errors.New("Unable to get data from match endpoint")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return matchHistory, errors.New("Unable to read response body")
	}

	errUnmarshal := json.Unmarshal(body, &matchHistory) 
	if errUnmarshal != nil {
		retrieveDataErr := fmt.Sprintf("Could not retrieve data for %s. Check if you linked blitz.gg with your account.", player.Name)
		return matchHistory, errors.New(retrieveDataErr)
	}

	return matchHistory, nil
}


func ParseValorantData(nametag string, blitzJson []byte) (definition.ValorantStats, error) {
	var valorantStats definition.ValorantStats
	err := json.Unmarshal(blitzJson, &valorantStats)
	if err != nil {
		retrieveDataErr := fmt.Sprintf("Could not retrieve data for %s. Check if you linked blitz.gg with your account.", nametag)
		return valorantStats, errors.New(retrieveDataErr)
	}
	return valorantStats, nil
}