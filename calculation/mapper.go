package calculation

func CreateCompetitiveTier(competitiveTier int) string {
	switch competitiveTier {
		case 0:
			return "Unranked"
		// TODO: what is up with these?
		// case 1:
		// case 2:
		case 3:
			return "Iron 1"
		case 4:
			return "Iron 2"
		case 5:
			return "Iron 3"
		case 6:
			return "Bronze 1"
		case 7:
			return "Bronze 2"
		case 8:
			return "Bronze 3"
		case 9:
			return "Silver 1"
		case 10:
			return "Silver 2"
		case 11:
			return "Silver 3"
		case 12:
			return "Gold 1"
		case 13:
			return "Gold 2"
		case 14:
			return "Gold 3"
		case 15:
			return "Platinum 1"
		case 16:
			return "Platinum 2"
		case 17:
			return "Platinum 3"
		case 18:
			return "Diamond 1"
		case 19:
			return "Diamond 2"
		case 20:
			return "Diamond 3"
		case 21:
			return "Immortal 1"
		case 22:
			return "Immortal 2"
		case 23:
			return "Immortal 3"
		case 24:
			return "Radiant"
	}
	return ""
}	