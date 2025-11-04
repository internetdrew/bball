package nba

type Team struct {
	ID                int    `json:"teamId"`
	Name              string `json:"teamName"`
	Tricode           string `json:"teamTricode"` // e.g. "NYK"
	City              string `json:"teamCity"`
	Score             int    `json:"score"`
	Wins              int    `json:"wins"`
	Losses            int    `json:"losses"`
	TimeoutsRemaining int    `json:"timeoutsRemaining,omitempty"`
}

type Leader struct {
	PersonID    int     `json:"personId"`
	Name        string  `json:"name"`
	JerseyNum   string  `json:"jerseyNum"`
	Position    string  `json:"position"`
	TeamTricode string  `json:"teamTricode"`
	PlayerSlug  *string `json:"playerSlug"`
	Points      int     `json:"points"`
	Rebounds    int     `json:"rebounds"`
	Assists     int     `json:"assists"`
}

type GameLeaders struct {
	HomeLeaders Leader `json:"homeLeaders"`
	AwayLeaders Leader `json:"awayLeaders"`
}

type Game struct {
	ID              string      `json:"gameId"`
	GameCode        string      `json:"gameCode"`
	GameStatus      int         `json:"gameStatus"`
	GameStatusText  string      `json:"gameStatusText"` // e.g. "In Progress", "Final", "Scheduled"
	Period          int         `json:"period,omitempty"`
	GameClock       string      `json:"gameClock,omitempty"`
	GameTimeUTC     string      `json:"gameTimeUTC"`
	GameDateTimeUTC string      `json:"gameDateTimeUTC"`
	GameDateTimeEst string      `json:"gameDateTimeEst"`
	GameEt          string      `json:"gameEt"`
	HomeTeam        Team        `json:"homeTeam"`
	AwayTeam        Team        `json:"awayTeam"`
	GameLeaders     GameLeaders `json:"gameLeaders,omitempty"`
}

type PlayerStats struct {
	PlayerName string `json:"player_name"`
	TeamCode   string `json:"team_code"`
	Points     int    `json:"points"`
	Rebounds   int    `json:"rebounds"`
	Assists    int    `json:"assists"`
}

type GameSummary struct {
	Game          Game          `json:"game"`
	TopPerformers []PlayerStats `json:"top_performers,omitempty"`
	LastUpdated   string        `json:"last_updated"`
}
