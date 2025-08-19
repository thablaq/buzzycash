package virtual

import (
"strings"
"errors"
)



type StartGameRequest struct {
	GameType string `json:"gameType" binding:"required,min=1"`
}

func (r *StartGameRequest) Validate() error{
	if strings.TrimSpace(r.GameType) == ""{
		return errors.New("gameType is required")
	}
	return nil 
}