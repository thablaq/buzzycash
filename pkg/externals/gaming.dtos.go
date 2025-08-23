package externals


// TokenResponse represents the response from the login endpoint
type TokenResponse struct {
	Accesstoken string `json:"Accesstoken"`
}

// RegisterUserRequest represents the request payload for user registration
type RegisterUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CompanyID string `json:"company_id"`
}

// GameRequest represents common game request structure
type GameRequest struct {
	GameID    string `json:"game_id"`
	CompanyID string `json:"company_id"`
}

// BuyTicketRequest represents the request payload for buying tickets
type BuyTicketRequest struct {
	GameID     string  `json:"game_id"`
	Username   string  `json:"username"`
	Quantity   int     `json:"quantity"`
	AmountPaid float64 `json:"amount_paid"`
	CompanyID  string  `json:"company_id"`
}

// CreateGameRequest represents the request payload for creating games
type CreateGameRequest struct {
	CompanyID            string  `json:"company_id"`
	GameName             string  `json:"game_name"`
	Amount               float64 `json:"amount"`
	DrawInterval         int     `json:"draw_interval"`
	WinningPercentage    float64 `json:"winning_percentage"`
	MaxWinners           int     `json:"max_winners"`
	Date                 string  `json:"date"`
	WeightedDistribution bool    `json:"weighted_distribution"`
}

// DebitWalletRequest represents the request payload for debiting wallet
type DebitWalletRequest struct {
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	CompanyID string  `json:"company_id"`
}

// PaymentRequest represents the request payload for payment
type CreditWalletRequest struct {
	UserID string  `json:"user_id"`
	Amount   float64 `json:"amount"`
	CompanyID string  `json:"company_id"`
}

// VirtualGameRequest represents the request payload for virtual games
type VirtualGameRequest struct {
	GameType string `json:"gameType"`
	Username string `json:"username"`
}


type BuyTicketResponse struct {
    TicketIDs []string `json:"ticket_ids"`
}

type PaymentResponse struct {
   Message string  `json:"message"`
}