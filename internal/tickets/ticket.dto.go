package tickets




type BuyTicketRequest struct {
	GameID     string  `json:"game_id" binding:"required" validate:"required"`
	Quantity   int     `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	AmountPaid float64 `json:"amount_paid" binding:"required,gt=0" validate:"required,gt=0"`
}

