package tickets




type BuyTicketRequest struct {
	GameID     string  `json:"game_id" binding:"required" validate:"required"`
	Quantity   int     `json:"quantity" binding:"required,gt=0" validate:"required,gt=0"`
	AmountPaid float64 `json:"amount_paid" binding:"required,gt=0" validate:"required,gt=0"`
}

type CreateGameRequest struct {
	GameName             string  `json:"game_name" binding:"required" validate:"required"`
	Amount               float64 `json:"amount" binding:"required" validate:"required"`
	DrawInterval         int     `json:"draw_interval" binding:"required" validate:"required"`
	WinningPercentage    float64 `json:"winning_percentage" binding:"required" validate:"required"`
	MaxWinners           int     `json:"max_winners" binding:"required" validate:"required"`
	Date                 string  `json:"date" binding:"required" validate:"required"`
	WeightedDistribution bool    `json:"weighted_distribution" validate:"required"`
}
