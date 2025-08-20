package externals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "net/url"
	"time"
	 "github.com/dblaq/buzzycash/internal/config"
)

// GamingService handles all gaming-related API operations
type GamingService struct {
	client *http.Client
}

// NewGamingService creates a new instance of GamingService
func NewGamingService() *GamingService {
	return &GamingService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}



// getAccessToken retrieves access token from the gaming API
func (gs *GamingService) getAccessToken() (string, error) {
	loginData := map[string]string{
		"username":   config.AppConfig.BuzzyCashUsername,
		"password":   config.AppConfig.BuzzyCashPassword,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println("Failed to marshal login data: ", err)
		return "", fmt.Errorf("failed to marshal login data: %w", err)
	}

	url := config.AppConfig.MaekandexGamingUrl + "login/"
	resp, err := gs.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to get access token: ", err)
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get access token, status: ", resp.StatusCode)
		return "", fmt.Errorf("failed to get token, status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		fmt.Println("Failed to decode token response: ", err)
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	fmt.Println("Access token retrieved successfully")
	return tokenResp.Accesstoken, nil
}

// makeAuthenticatedRequest makes an HTTP request with authentication
func (gs *GamingService) makeAuthenticatedRequest(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	token, err := gs.getAccessToken()
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	baseURL := config.AppConfig.MaekandexGamingUrl + endpoint
	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Add query parameters if provided
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return gs.client.Do(req)
}

// RegisterUser registers a new user in the gaming system
func (gs *GamingService) RegisterUser(username, email, firstName, lastName string) (map[string]interface{}, error) {
	reqData := RegisterUserRequest{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "register/user/", reqData, nil)
	if err != nil {
		fmt.Println("Error registering user: ", err)
		return nil, fmt.Errorf("failed to register profile: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Profile api responded")
	return result, nil
}

// StartGame starts a game
func (gs *GamingService) StartGame(gameID string) (map[string]interface{}, error) {
	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "games/start/", reqData, nil)
	if err != nil {
		fmt.Println("Failed to start game: ", err)
		return nil, fmt.Errorf("failed to start game: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Game started successfully: ", result)
	return result, nil
}

// StopGame stops a game
func (gs *GamingService) StopGame(gameID string) (map[string]interface{}, error) {
	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "games/stop/", reqData, nil)
	if err != nil {
		fmt.Println("Failed to stop game: ", err)
		return nil, fmt.Errorf("failed to stop game: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Game stopped successfully: ", result)
	return result, nil
}

// GetDraws retrieves draws for a game
func (gs *GamingService) GetDraws(gameID string) (map[string]interface{}, error) {
	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "games/draw/", reqData, nil)
	if err != nil {
		fmt.Println("Failed to get draws: ", err)
		return nil, fmt.Errorf("failed to get draws: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Draws retrieved successfully: ", result)
	return result, nil
}

// GetAllTickets retrieves all tickets
func (gs *GamingService) GetAllTickets() (map[string]interface{}, error) {
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "admin/get_tickets/", nil, params)
	if err != nil {
		fmt.Println("Failed to get tickets: ", err)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Tickets retrieved successfully: ", result)
	return result, nil
}

// GetWalletBalances retrieves all wallet balances
func (gs *GamingService) GetWalletBalances() (map[string]interface{}, error) {
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "all/wallet", nil, params)
	if err != nil {
		fmt.Println("Failed to get user balances: ", err)
		return nil, fmt.Errorf("failed to get user balances: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("User balances retrieved successfully: ", result)
	return result, nil
}

// GetUserWallet retrieves wallet for a specific user
func (gs *GamingService) GetUserWallet(username string) (map[string]interface{}, error) {
	params := map[string]string{
		"user_id": username,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "check/wallet", nil, params)
	if err != nil {
		fmt.Println("Failed to get wallet for username: ", err)
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("User wallet retrieved successfully for username")
	return result, nil
}

// BuyTicket purchases a ticket for a game
func (gs *GamingService) BuyTicket(gameID, username string, quantity int, amountPaid float64) (*BuyTicketResponse, error) {
	reqData := BuyTicketRequest{
		GameID:     gameID,
		Username:   username,
		Quantity:   quantity,
		AmountPaid: amountPaid,
		CompanyID:  config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "games/buy/", reqData, nil)
	if err != nil {
		fmt.Println("Failed to purchase ticket: ", err)
		return nil, fmt.Errorf("failed to purchase ticket: %w", err)
	}
	defer resp.Body.Close()

	var result BuyTicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Ticket purchased successfully: ", result)
	return &result, nil
}

// GetUserTickets retrieves tickets for a specific user
func (gs *GamingService) GetUserTickets(username string) (map[string]interface{}, error) {
	params := map[string]string{
		"user_id":    username,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "users/tickets", nil, params)
	if err != nil {
		fmt.Println("Failed to get tickets for username: ", err)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("User tickets retrieved successfully for username")
	return result, nil
}

// GetUserResults retrieves results for a specific user
func (gs *GamingService) GetUserResults(username string) (map[string]interface{}, error) {
	params := map[string]string{
		"username":   username,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "games/results", nil, params)
	if err != nil {
		fmt.Println("Failed to get results for username: ", err)
		return nil, fmt.Errorf("failed to get results: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("User results retrieved successfully for username")
	return result, nil
}

// GetWinnerLogs retrieves winner logs
func (gs *GamingService) GetWinnerLogs() (map[string]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "winners/logs/", nil, nil)
	if err != nil {
		fmt.Println("Failed to get winner logs: ", err)
		return nil, fmt.Errorf("failed to get winner logs: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Winner logs retrieved successfully: ", result)
	return result, nil
}

// GetLeaderBoard retrieves the leaderboard
func (gs *GamingService) GetLeaderBoard() (map[string]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "leaderboard", nil, nil)
	if err != nil {
		fmt.Println("Failed to get leaderboard: ", err)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Leaderboard retrieved successfully: ", result)
	return result, nil
}

// GetVirtualGames retrieves available virtual games
func (gs *GamingService) GetVirtualGames() ([]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "list/virtual/games", nil, nil)
	if err != nil {
		fmt.Println("Failed to get virtual games: ", err)
		return nil, fmt.Errorf("failed to get virtual games: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	data, ok := result["data"]
	if !ok {
		fmt.Println("Invalid games data received from Meakindex")
		return nil, fmt.Errorf("invalid games data received from Meakindex")
	}

	gamesList, ok := data.([]interface{})
	if !ok {
		fmt.Println("Invalid games data format received from Meakindex")
		return nil, fmt.Errorf("invalid games data format received from Meakindex")
	}

	fmt.Println(gamesList)
	fmt.Println("Virtual games retrieved successfully")
	return gamesList, nil
}

// StartVirtualGame starts a virtual game
func (gs *GamingService) StartVirtualGame(gameType, username string) (map[string]interface{}, error) {
	params := map[string]string{
		"gameType": gameType,
		"username": username,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "virtualstart/game", nil, params)
	if err != nil {
		fmt.Println("Failed to start virtual game for gameType: ", err)
		return nil, fmt.Errorf("failed to start virtual game: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println(result)
	fmt.Println("Virtual game started successfully for gameType")
	return result, nil
}

// CreateGames creates a new game (admin function)
func (gs *GamingService) CreateGames(gameName string, amount float64, drawInterval int, winningPercentage float64, maxWinners int, date string, weightedDistribution bool) (map[string]interface{}, error) {
	reqData := CreateGameRequest{
		CompanyID:            config.AppConfig.BuzzyCashCompanyID,
		GameName:             gameName,
		Amount:               amount,
		DrawInterval:         drawInterval,
		WinningPercentage:    winningPercentage,
		MaxWinners:           maxWinners,
		Date:                 date,
		WeightedDistribution: weightedDistribution,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "create/games/", reqData, nil)
	if err != nil {
		fmt.Println("Error creating game: ", err)
		return nil, fmt.Errorf("failed to create game with Meakindex: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetGames retrieves all games
func (gs *GamingService) GetGames() (map[string]interface{}, error) {
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "games/", nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get games: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// DebitUserWallet debits amount from user's wallet
func (gs *GamingService) DebitUserWallet(phoneNumber string, amount float64) (map[string]interface{}, error) {
	reqData := DebitWalletRequest{
		UserID:    phoneNumber,
		Amount:    amount,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "debit/wallet/", reqData, nil)
	if err != nil {
		fmt.Println("Failed to debit wallet: ", err)
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Wallet debited successfully: ", result)
	return result, nil
}

// GetPaymentLink gets payment link for user
func (gs *GamingService) GetPaymentLink(username string, amount float64) (*PaymentLinkResponse, error) {
	reqData := PaymentRequest{
		Username: username,
		Amount:   amount,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "payment/", reqData, nil)
	if err != nil {
		fmt.Println("Error getting payment link: ", err)
		return nil, fmt.Errorf("failed to get payment link: %w", err)
	}
	defer resp.Body.Close()

	var result PaymentLinkResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// VerifyPayment verifies a payment
func (gs *GamingService) VerifyPayment(paymentID string) (map[string]interface{}, error) {
	params := map[string]string{
		"payment_id": paymentID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "verify/payment/", nil, params)
	if err != nil {
		fmt.Println("Error verifying payment: ", err)
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ListPayouts lists all payouts
func (gs *GamingService) ListPayouts() (map[string]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "payouts/", nil, nil)
	if err != nil {
		fmt.Println("Error listing payouts: ", err)
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ListUserPayout lists payouts for a specific user
func (gs *GamingService) ListUserPayout(username string) (map[string]interface{}, error) {
	params := map[string]string{
		"username": username,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "payouts/", nil, params)
	if err != nil {
		fmt.Println("Error listing user payouts: ", err)
		return nil, fmt.Errorf("failed to list user payouts: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// PayoutByAdmin processes payout by admin
func (gs *GamingService) PayoutByAdmin(payoutID string) (map[string]interface{}, error) {
	params := map[string]string{
		"payout_id": payoutID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "payouts/", nil, params)
	if err != nil {
		fmt.Println("Error processing payout: ", err)
		return nil, fmt.Errorf("failed to process payout: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}