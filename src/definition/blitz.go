package definition

import (
	"time"
)

type Stats struct {
	Overall Overall
}

type Ranks struct {
	Competitive struct {
		Tier int 
	} `json:"competitive"`
}

type ValorantStats struct {
	Name string
	Id string
	Stats Stats	
	Ranks Ranks
}

type Overall struct {
	Career MatchStatistics
	Last20 MatchStatistics
}

type RoundResult struct { 
	RoundNum    int    	`json:"roundNum"`
	PlantSite   string `json:"plantSite"`
	BombPlanter string `json:"bombPlanter,omitempty"`
	PlayerStats []PlayerStat `json:"playerStats"`
	RoundResult  string `json:"roundResult"`
	WinningTeam  string `json:"winningTeam"`
}

type Player struct {
	Stats  					PlayerStats	`json:"stats"`
	TeamID      			string 		`json:"teamId"`			// blue or red
	PartyID    	 			string 		`json:"partyId"`
	Subject     			string 		`json:"subject"`			// player id
	CharacterID 			string 		`json:"characterId"`		// agent
	CompetitiveTier         int 		`json:"competitiveTier"`
	SessionPlaytimeMinutes  int 		`json:"sessionPlaytimeMinutes"`
	RoundDamage []struct {
		Round    int    `json:"round"`
		Damage   int    `json:"damage"`
		Receiver string `json:"receiver"`
	} `json:"roundDamage"`
}

type DamageStats struct {
	Bodyshots int
	Headshots int
	Legshots  int
	Damage    int
}

type MatchHistoryOffset struct {
	Count int
	Data []MatchHistory
	Limit int
	Offset string
}

type MatchHistory struct {
	ID     		 string 	   `json:"id"`
	Map    	     string 	   `json:"map"`
	Mode   		 string 	   `json:"mode"`
	Ranked 		 bool   	   `json:"ranked"`
	Teams  		 []Team 	   `json:"teams"`
	RoundResults []RoundResult `json:"roundResults"`
	StartedAt 	 time.Time 	   `json:"startedAt"`
	Players 	 []Player      `json:"players"`
	Length    	 int       	   `json:"length"`
	Queue        string        `json:"queue"`
	Season    	 string        `json:"season"`
	Version   	 string        `json:"version"`
}

type Team struct {
	Won          bool   `json:"won"`
	TeamID       string `json:"teamId"`
	RoundsWon    int    `json:"roundsWon"`
	RoundsPlayed int    `json:"roundsPlayed"`
}

type MatchStatistics struct {
	DamageStats DamageStats
	Matches int
}

type Damage struct {
	Receiver  string `json:"receiver"`
	Damage    int    `json:"damage"`
	Legshots  int    `json:"legshots"`
	Bodyshots int    `json:"bodyshots"`
	Headshots int    `json:"headshots"`
}

type PlayerStat struct {
	Score  int `json:"score"`
	Damage []Damage `json:"damage"`
	WasAfk  bool `json:"wasAfk"`
	Subject      string `json:"subject"`
	WasPenalized bool   `json:"wasPenalized"`
}

type PlayerStats struct {
	Kills        int `json:"kills"`
	Score        int `json:"score"`
	Deaths       int `json:"deaths"`
	Assists      int `json:"assists"`
	AbilityCasts struct {
		GrenadeCasts  int `json:"grenadeCasts"`
		Ability1Casts int `json:"ability1Casts"`
		Ability2Casts int `json:"ability2Casts"`
		UltimateCasts int `json:"ultimateCasts"`
	} `json:"abilityCasts"`
	RoundsPlayed   int `json:"roundsPlayed"`
	PlaytimeMillis int `json:"playtimeMillis"`
}
