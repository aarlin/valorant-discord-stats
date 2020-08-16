package formatter

import (
	"github.com/aarlin/valorant-discord-stats/definition"
	"github.com/aarlin/valorant-discord-stats/calculation"
	"fmt"
	"errors"
)

func GenerateMatchSummary(player definition.ValorantStats, matchHistory definition.MatchHistory) (string, error) {
	switch matchHistory.Queue {
		case "competitive": return generateRegularMatchSummary(player, matchHistory)
		case "deathmatch": return generateDeathMatchSummary(player, matchHistory)
	}
	return "", errors.New("Unable to determine match queue type")
}

func generateHitPercentages(player definition.ValorantStats, matchHistory definition.MatchHistory) map[string]*definition.DamageStats {
	var matchDamageStatistics = make(map[string]*definition.DamageStats)
	// make call to each match id using get request 
	// then do this nested for loop to grab all damage
	matchDamageStatistics[matchHistory.ID] = &definition.DamageStats{0, 0, 0, 0}
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
	return matchDamageStatistics
}

func generateGameRoundResults(playerTeam string, matchHistory definition.MatchHistory) string {
	// TODO: parse round damage for dmg/round
	// var roundDamage struct
	gameRoundResults := "0 - 0"

	var blueTeamRoundWins int
	var redTeamRoundWins int

	for _, matchParticipantTeam := range matchHistory.Teams {
		switch matchParticipantTeam.TeamID {
			case "Red": redTeamRoundWins = matchParticipantTeam.RoundsWon
			case "Blue": blueTeamRoundWins = matchParticipantTeam.RoundsWon
			default: fmt.Sprintf("Something went wrong parsing team ID from API endpoint")
		}
	}

	switch playerTeam {
		case "Red":	gameRoundResults = fmt.Sprintf("%d - %d", redTeamRoundWins, blueTeamRoundWins)
		case "Blue": gameRoundResults = fmt.Sprintf("%d - %d", blueTeamRoundWins, redTeamRoundWins)
		default: fmt.Sprintf("Something went wrong while trying to create round results")
	}
	
	return gameRoundResults
}

func generateRegularMatchSummary(player definition.ValorantStats, matchHistory definition.MatchHistory) (string, error) {
	var kills int
	var deaths int
	var assists int
	var competitiveTier int
	var score int
	var roundsPlayed int
	var team string 
	var gameRoundResults string

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

	gameRoundResults = generateGameRoundResults(team, matchHistory)

	var matchPercentages = calculation.CalculateHitPercentages(*matchDamageStatistics[matchHistory.ID])
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

func generateDeathMatchSummary(player definition.ValorantStats, matchHistory definition.MatchHistory) (string, error) {

	return "", nil
}

func GenerateDiscordEmbedContent(nametag string, stats definition.ValorantStats, hitRate definition.HitPercentages, matchStatisticType string) string {
	headShots := fmt.Sprintf(":no_mouth: Head shot percentage: %.2f%%\n", hitRate.HeadShotPercentage)
	bodyShots := fmt.Sprintf(":shirt: Body shot percentage: %.2f%%\n", hitRate.BodyShotPercentage)
	legShots := fmt.Sprintf(":foot: Leg shot percentage: %.2f%%\n", hitRate.LegShotPercentage)
	matchesPlayed := stats.Stats.Overall.Career.Matches
	content := ""
	switch matchStatisticType {
	case "career":
		content = fmt.Sprintf("Career Stats for %s:\nTotal number of matches: %d\n%s%s%s\n", nametag, matchesPlayed, headShots, bodyShots, legShots)
	case "last20":
		content = fmt.Sprintf("Last 20 Games Stats for %s:\n%s%s%s\n", nametag, headShots, bodyShots, legShots)
	default:
		content = fmt.Sprintf("Something went wrong choosing the statistic type when posting to Discord")
	}

	return content
}
