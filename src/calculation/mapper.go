package calculation

import (
	"github.com/aarlin/valorant-discord-stats/src/config"
)

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


func GenerateCompetitiveTierImage(competitiveTier int) string {
	switch competitiveTier {
		case 0:
			return config.UNRANKED_LINK
		// TODO: what is up with these?
		// case 1:
		// case 2:
		case 3:
			return config.IRON_1_LINK
		case 4:
			return config.IRON_2_LINK
		case 5:
			return config.IRON_3_LINK
		case 6:
			return config.BRONZE_1_LINK
		case 7:
			return config.BRONZE_2_LINK
		case 8:
			return config.BRONZE_3_LINK
		case 9:
			return config.SILVER_1_LINK
		case 10:
			return config.SILVER_2_LINK
		case 11:
			return config.SILVER_3_LINK
		case 12:
			return config.GOLD_1_LINK
		case 13:
			return config.GOLD_2_LINK
		case 14:
			return config.GOLD_3_LINK
		case 15:
			return config.PLATINUM_1_LINK
		case 16:
			return config.PLATINUM_2_LINK
		case 17:
			return config.PLATINUM_3_LINK
		case 18:
			return config.DIAMOND_1_LINK
		case 19:
			return config.DIAMOND_2_LINK
		case 20:
			return config.DIAMOND_3_LINK
		case 21:
			return config.IMMORTAL_1_LINK
		case 22:
			return config.IMMORTAL_2_LINK
		case 23:
			return config.IMMORTAL_3_LINK
		case 24:
			return config.RADIANT_LINK
	}
	return ""
}