package formatter

import (
	"github.com/aarlin/valorant-discord-stats/definition"
	"github.com/aarlin/valorant-discord-stats/calculation"
	"fmt"
	"errors"
)

func GenerateMatchSummaryText(player definition.ValorantStats, matchHistory definition.MatchHistory) string {
	switch matchHistory.Queue {
		case "competitive": 
			return generateRegularMatchSummaryText(matchSummary)
		case "deathmatch": 
			return generateDeathMatchSummary(player, matchHistory)
	}
	return "", errors.New("Unable to determine match queue type")
}

func generateHitCount(player definition.ValorantStats, matchHistory definition.MatchHistory) map[string]*definition.DamageStats {
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

func generatePlayerDamagePerRound(player definition.ValorantStats, matchHistory definition.MatchHistory) map[int]int {
	var playerDamagePerRound = make(map[int]int)
	for _, matchParticipant := range matchHistory.Players {
		if matchParticipant.Subject == player.Id {
			for _, round := range matchParticipant.RoundDamage {
				if val, ok := playerDamagePerRound[round.Round]; ok {
					// mapping is correct BUT missing rounds where player does no damage
					playerDamagePerRound[round.Round] = val + round.Damage
				} else {
					playerDamagePerRound[round.Round] = round.Damage
				}
			}
		}
	}
	return playerDamagePerRound
}

func generateRegularMatchSummary(player definition.ValorantStats, matchHistory definition.MatchHistory) definition.RegularMatchSummary {
	var kills int
	var deaths int
	var assists int
	var competitiveTier int
	var score int
	var roundsPlayed int
	var team string 
	var gameRoundResults string
	var hitCount = make(map[string]*definition.DamageStats)
	var playerDamagePerRound = make(map[int]int)

	for _, matchParticipant := range matchHistory.Players {
		if matchParticipant.Subject == player.Id {
			kills = matchParticipant.Stats.Kills
			deaths = matchParticipant.Stats.Deaths
			assists = matchParticipant.Stats.Assists
			competitiveTier = matchParticipant.CompetitiveTier
			score = matchParticipant.Stats.Score
			roundsPlayed = matchParticipant.Stats.RoundsPlayed
			team = matchParticipant.TeamID
		}
	}

	playerDamagePerRound = generatePlayerDamagePerRound(player, matchHistory)
	gameRoundResults = generateGameRoundResults(team, matchHistory)
	hitCount = generateHitCount(player, matchHistory)

	var matchPercentages = calculation.CalculateHitPercentages(*hitCount[matchHistory.ID])

	var regularMatchSummary = definition.RegularMatchSummary {
		Nametag: player.Nametag,
		CompetitiveTier: calculation.CreateCompetitiveTier(competitiveTier),
		GameRoundResults: gameRoundResults,
		MatchHistoryID: matchHistory.ID,
		MatchHistoryMap: matchHistory.Map,
		Headshots: hitCount[matchHistory.ID].Headshots,
		HeadShotPercentage: matchPercentages.HeadShotPercentage,
		Bodyshots: hitCount[matchHistory.ID].Bodyshots,
		BodyShotPercentage: matchPercentages.BodyShotPercentage,
		Legshots: hitCount[matchHistory.ID].Legshots,
		LegShotPercentage: matchPercentages.LegShotPercentage,
		Damage: hitCount[matchHistory.ID].Damage,
		CombatScore: (score / roundsPlayed),
		Kills: kills,
		Deaths: deaths,
		Assists: assists,
	}
	
	return regularMatchSummary
}

func generateRegularMatchSummaryText(matchSummary definition.RegularMatchSummary) string {
	// TODO: remove map
	// Remove rank
	// Add those into images
	// Turn map ID into link as a sub title
	
	// TODO: Add map image 
	// TODO: Add competitive tier image
	// var mapImage = fmt.Sprintf("https://blitz-cdn.blitz.gg/blitz/val/maps/map-art-%s.jpg", matchHistory.Map)

	var matchStats = fmt.Sprintf("Nametag: %s\n" + 
		"Game Results: %s\n" + 
		"Headshots: %d (%.2f%%)\n" +
		"Bodyshots: %d (%.2f%%)\n" + 
		"Legshots: %d (%.2f%%)\n" + 
		"Damage: %d\n" + 
		"Combat Score: %d\n" + 
		"K\\/D\\/A: %d\\/%d\\/%d\n",
		matchSummary.Nametag,
		matchSummary.GameRoundResults,
		matchSummary.Headshots, matchSummary.HeadShotPercentage,
		matchSummary.Bodyshots, matchSummary.BodyShotPercentage,
		matchSummary.Legshots, matchSummary.LegShotPercentage,
		matchSummary.Damage,
		matchSummary.CombatScore,
		matchSummary.Kills, 
		matchSummary.Deaths, 
		matchSummary.Assists)

	return matchStats
}

func generateDeathMatchSummary(player definition.ValorantStats, matchHistory definition.MatchHistory) definition.DeathMatchSummary {
	var matchSummary = definition.DeathMatchSummary {

	}
	return matchSummary
}

func generateDeathMatchSummaryText(matchSummary definition.DeathMatchSummary) string {
	return ""
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
