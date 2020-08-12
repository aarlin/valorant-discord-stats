package handler

import (
	"math"
	"github.com/dariubs/percent"
	"github.com/aarlin/valorant-discord-stats/definition"
)

func calculateHitPercentages(damageStats definition.DamageStats) definition.HitPercentages {
	headShots := damageStats.Headshots
	bodyShots := damageStats.Bodyshots
	legShots := damageStats.Legshots

	var hitPercentages definition.HitPercentages
	totalShots := headShots + bodyShots + legShots
	hitPercentages.HeadShotPercentage = roundPercentage(percent.PercentOf(headShots, totalShots))
	hitPercentages.BodyShotPercentage = roundPercentage(percent.PercentOf(bodyShots, totalShots))
	hitPercentages.LegShotPercentage  = roundPercentage(percent.PercentOf(legShots, totalShots))
	
	return hitPercentages
}

func roundPercentage(percentage float64) float64 {
	return math.Round(percentage * 100) / 100
}