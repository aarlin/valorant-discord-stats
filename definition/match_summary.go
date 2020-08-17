package definition

type RegularMatchSummary struct {
	Nametag string
	CompetitiveTier string
	GameRoundResults string
	MatchHistoryID string
	MatchHistoryMap string
	Headshots int
	HeadShotPercentage float64
	Bodyshots int
	BodyShotPercentage float64
	Legshots int
	LegShotPercentage float64
	Damage int
	CombatScore int
	Kills int
	Deaths int
	Assists int
}

type DeathMatchSummary struct {
	Placement int
}

type MatchSummary struct {
	Kills int
	Deaths int
	Assists int
	CompetitiveTier int
	Score int
	RoundsPlayed int
	Team string
	Queue string
	RegularMatch RegularMatchSummary
	DeathMatch DeathMatchSummary
	
}