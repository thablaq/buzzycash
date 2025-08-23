package virtual

import (
"strings"
"errors"
)



type StartGameRequest struct {
	GameType string `json:"game_type" binding:"required,min=1"`
}

func (r *StartGameRequest) Validate() error{
	if strings.TrimSpace(r.GameType) == ""{
		return errors.New("game_type is required")
	}
	return nil 
}