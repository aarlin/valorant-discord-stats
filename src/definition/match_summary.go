package definition

type MatchSummary struct {
	Name string
	Kills int
	Deaths int
	Assists int
	CompetitiveTier int
	Rank string
	Score int
	RoundsPlayed int
	Team string
	Queue string
	ID string
	Map string
	RegularMatch RegularMatchSummary
	DeathMatch DeathMatchSummary
}

type RegularMatchSummary struct {
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
}

type DeathMatchSummary struct {
	Placement int
}