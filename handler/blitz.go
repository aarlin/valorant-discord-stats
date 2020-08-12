package handler

import (
	"fmt"
	"net/http"
	"errors"
	"io/ioutil"
	"encoding/json"
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

func RetrieveMatchStats(player ValorantStats, matches []string) (string, error) {
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