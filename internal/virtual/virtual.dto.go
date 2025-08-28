package virtual




type StartGameRequest struct {
	GameType string `json:"game_type" binding:"required"`
}

