package calculation

func CreateCompetitiveTier(competitiveTier int) string {
	switch competitiveTier {
		case 0:
			return "unranked"
		// TODO: what is up with these?
		// case 1:
		// case 2:
		case 3:
			return "iron1"
		case 4:
			return "iron2"
		case 5:
			return "iron3"
		case 6:
			return "bronze1"
		case 7:
			return "bronze2"
		case 8:
			return "bronze3"
		case 9:
			return "silver1"
		case 10:
			return "silver2"
		case 11:
			return "silver3"
		case 12:
			return "gold1"
		case 13:
			return "gold2"
		case 14:
			return "gold3"
		case 15:
			return "platinum1"
		case 16:
			return "platinum2"
		case 17:
			return "platinum3"
		case 18:
			return "diamond1"
		case 19:
			return "diamond2"
		case 20:
			return "diamond3"
		case 21:
			return "immortal1"
		case 22:
			return "immortal2"
		case 23:
			return "immortal3"
		case 24:
			return "radiant"
	}
	// TODO: Add rank icon https://blitz-cdn-plain.blitz.gg/blitz/val/ranks/diamond_small.svg
	// https://blitz-cdn-plain.blitz.gg/blitz/val/ranks/gold1.svg
	// Create rank average in match
	return ""
}	