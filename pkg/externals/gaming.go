package externals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
     "log"
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
	log.Println("Starting getAccessToken process")

	loginData := map[string]string{
		"username":   config.AppConfig.BuzzyCashUsername,
		"password":   config.AppConfig.BuzzyCashPassword,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Login data prepared")

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		log.Println("Failed to marshal login data: ", err)
		return "", fmt.Errorf("failed to marshal login data: %w", err)
	}
	log.Println("Login data marshaled successfully")

	url := config.AppConfig.MaekandexGamingUrl + "login/"
	log.Println("Sending POST request to URL:")

	resp, err := gs.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("Failed to get access token: ", err)
		return "", fmt.Errorf("failed to get token: %w", err)
	}
	defer resp.Body.Close()
	log.Println("POST request sent successfully, awaiting response")

	if resp.StatusCode != http.StatusOK {
		log.Println("Failed to get access token, status: ", resp.StatusCode)
		return "", fmt.Errorf("failed to get token, status: %d", resp.StatusCode)
	}
	log.Println("Response received with status code:", resp.StatusCode)

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		log.Println("Failed to decode token response: ", err)
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}
	log.Println("Token response decoded successfully")

	log.Println("Access token retrieved successfully:")
	return tokenResp.Accesstoken, nil
}

// makeAuthenticatedRequest makes an HTTP request with authentication
func (gs *GamingService) makeAuthenticatedRequest(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	log.Println("Starting makeAuthenticatedRequest")

	token, err := gs.getAccessToken()
	if err != nil {
		log.Println("Failed to get access token: ",err)
		return nil, err
	}
	log.Println("Access token retrieved successfully")

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			log.Println("Failed to marshal request body: ", err)
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
		log.Println("Request body marshaled successfully")
	}

	baseURL := config.AppConfig.MaekandexGamingUrl + endpoint
	req, err := http.NewRequest(method, baseURL, reqBody)
	if err != nil {
		log.Println("Failed to create request: ", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	log.Println("HTTP request created successfully")

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	log.Println("Headers set successfully")

	// Add query parameters if provided
	if params != nil {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
		log.Println("Query parameters added successfully")
	}

	log.Println("Sending HTTP request")
	resp, err := gs.client.Do(req)
	if err != nil {
		log.Println("Failed to send HTTP request: ", err)
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	log.Println("HTTP request sent successfully")
	return resp, nil
}

// RegisterUser registers a new user in the gaming system
func (gs *GamingService) RegisterUser(username, email, firstName, lastName string) (map[string]interface{}, error) {
	log.Println("Starting RegisterUser")

	reqData := RegisterUserRequest{
		Username:  username,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Request data prepared for RegisterUser")

	resp, err := gs.makeAuthenticatedRequest("POST", "register/user/", reqData, nil)
	if err != nil {
		log.Println("Error registering user: ", err)
		return nil, fmt.Errorf("failed to register profile: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for RegisterUser")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode RegisterUser response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("RegisterUser response decoded successfully")
	return result, nil
}

// StartGame starts a game
func (gs *GamingService) StartGame(gameID string) (map[string]interface{}, error) {
	log.Println("Starting StartGame")

	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Request data prepared for StartGame")

	resp, err := gs.makeAuthenticatedRequest("POST", "games/start/", reqData, nil)
	if err != nil {
		log.Println("Failed to start game: ", err)
		return nil, fmt.Errorf("failed to start game: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for StartGame")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode StartGame response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Game started successfully: ", result)
	return result, nil
}

// StopGame stops a game
func (gs *GamingService) StopGame(gameID string) (map[string]interface{}, error) {
	log.Println("Starting StopGame")

	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Request data prepared for StopGame")

	resp, err := gs.makeAuthenticatedRequest("POST", "games/stop/", reqData, nil)
	if err != nil {
		log.Println("Failed to stop game: ", err)
		return nil, fmt.Errorf("failed to stop game: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for StopGame")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode StopGame response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Game stopped successfully: ", result)
	return result, nil
}

// GetDraws retrieves draws for a game
func (gs *GamingService) GetDraws(gameID string) (map[string]interface{}, error) {
	log.Println("Starting GetDraws")

	reqData := GameRequest{
		GameID:    gameID,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Request data prepared for GetDraws")

	resp, err := gs.makeAuthenticatedRequest("POST", "games/draw/", reqData, nil)
	if err != nil {
		log.Println("Failed to get draws: ", err)
		return nil, fmt.Errorf("failed to get draws: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetDraws")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode GetDraws response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Draws retrieved successfully: ", result)
	return result, nil
}

// GetAllTickets retrieves all tickets
func (gs *GamingService) GetAllTickets() (map[string]interface{}, error) {
	log.Println("Starting GetAllTickets")

	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Query parameters prepared for GetAllTickets")

	resp, err := gs.makeAuthenticatedRequest("GET", "admin/get_tickets/", nil, params)
	if err != nil {
		log.Println("Failed to get tickets: ", err)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetAllTickets")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode GetAllTickets response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Tickets retrieved successfully: ", result)
	return result, nil
}

// GetWalletBalances retrieves all wallet balances
func (gs *GamingService) GetWalletBalances() (map[string]interface{}, error) {
	log.Println("Starting GetWalletBalances")

	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Query parameters prepared for GetWalletBalances")

	resp, err := gs.makeAuthenticatedRequest("GET", "all/wallet", nil, params)
	if err != nil {
		log.Println("Failed to get user balances: ", err)
		return nil, fmt.Errorf("failed to get user balances: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetWalletBalances")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode GetWalletBalances response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User balances retrieved successfully: ", result)
	return result, nil
}

// GetUserWallet retrieves wallet for a specific user
func (gs *GamingService) GetUserWallet(username string) (map[string]interface{}, error) {
	log.Println("Starting GetUserWallet")

	params := map[string]string{
		"user_id": username,
	}
	log.Println("Query parameters prepared for GetUserWallet")

	resp, err := gs.makeAuthenticatedRequest("GET", "check/wallet", nil, params)
	if err != nil {
		log.Println("Failed to get wallet for username: ", err)
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetUserWallet")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode GetUserWallet response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User wallet retrieved successfully for username")
	return result, nil
}

// BuyTicket purchases a ticket for a game
func (gs *GamingService) BuyTicket(gameID, username string, quantity int, amountPaid float64) (*BuyTicketResponse, error) {
	log.Println("Starting BuyTicket")

	reqData := BuyTicketRequest{
		GameID:     gameID,
		Username:   username,
		Quantity:   quantity,
		AmountPaid: amountPaid,
		CompanyID:  config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Request data prepared for BuyTicket")

	resp, err := gs.makeAuthenticatedRequest("POST", "games/buy/", reqData, nil)
	if err != nil {
		log.Println("Failed to purchase ticket: ", err)
		return nil, fmt.Errorf("failed to purchase ticket: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for BuyTicket")

	var result BuyTicketResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode BuyTicket response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("Ticket purchased successfully: ", result)
	return &result, nil
}

// GetUserTickets retrieves tickets for a specific user
func (gs *GamingService) GetUserTickets(username string) (map[string]interface{}, error) {
	log.Println("Starting GetUserTickets")

	params := map[string]string{
		"user_id":    username,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Query parameters prepared for GetUserTickets")

	resp, err := gs.makeAuthenticatedRequest("GET", "users/tickets", nil, params)
	if err != nil {
		log.Println("Failed to get tickets for username: ", err)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetUserTickets")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode GetUserTickets response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User tickets retrieved successfully for username")
	return result, nil
}

// GetUserResults retrieves results for a specific user
func (gs *GamingService) GetUserResults(username string) (map[string]interface{}, error) {
	log.Println("Starting GetUserResults")

	params := map[string]string{
		"username":   username,
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}
	log.Println("Query parameters prepared for GetUserResults")

	resp, err := gs.makeAuthenticatedRequest("GET", "games/results", nil, params)
	if err != nil {
		log.Println("Failed to get results for username: ", err)
		return nil, fmt.Errorf("failed to get results: %w", err)
	}
	defer resp.Body.Close()
	log.Println("Response received for GetUserResults")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode GetUserResults response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Println("User results retrieved successfully for username")
	return result, nil
}

// GetWinnerLogs retrieves winner logs

func (gs *GamingService) GetWinnerLogs() (map[string]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "winners/logs/", nil, nil)
	if err != nil {
		log.Println("Failed to get winner logs: ", err)
		return nil, fmt.Errorf("failed to get winner logs: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode winner logs response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Winner logs retrieved successfully: ", result)
	return result, nil
}

// GetLeaderBoard retrieves the leaderboard
func (gs *GamingService) GetLeaderBoard() (map[string]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "leaderboard", nil, nil)
	if err != nil {
		log.Println("Failed to get leaderboard: ", err)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode leaderboard response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Leaderboard retrieved successfully: ", result)
	return result, nil
}

// GetVirtualGames retrieves available virtual games
func (gs *GamingService) GetVirtualGames() ([]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "list/virtual/games", nil, nil)
	if err != nil {
		log.Println("Failed to get virtual games: ", err)
		return nil, fmt.Errorf("failed to get virtual games: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode virtual games response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	data, ok := result["data"]
	if !ok {
		log.Println("Invalid games data received from Meakindex")
		return nil, fmt.Errorf("invalid games data received from Meakindex")
	}

	gamesList, ok := data.([]interface{})
	if !ok {
		log.Println("Invalid games data format received from Meakindex")
		return nil, fmt.Errorf("invalid games data format received from Meakindex")
	}

	log.Println("Virtual games retrieved successfully: ", gamesList)
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
		log.Println("Failed to start virtual game for gameType: ", err)
		return nil, fmt.Errorf("failed to start virtual game: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode virtual game start response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Virtual game started successfully for gameType: ", result)
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
		log.Println("Error creating game: ", err)
		return nil, fmt.Errorf("failed to create game with Meakindex: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode create game response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Game created successfully: ", result)
	return result, nil
}

// GetGames retrieves all games
func (gs *GamingService) GetGames() (map[string]interface{}, error) {
	params := map[string]string{
		"company_id": config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "games/", nil, params)
	if err != nil {
		log.Println("Failed to get games: ", err)
		return nil, fmt.Errorf("failed to get games: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode games response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Games retrieved successfully: ", result)
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
		log.Println("Failed to debit wallet: ", err)
		return nil, fmt.Errorf("failed to debit wallet: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode debit wallet response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Wallet debited successfully: ", result)
	return result, nil
}

// GetPaymentLink gets payment link for user
func (gs *GamingService) CreditUserWallet(username string, amount float64) (*PaymentResponse, error) {
	reqData := CreditWalletRequest{
		UserID: username,
		Amount: amount,
		CompanyID: config.AppConfig.BuzzyCashCompanyID,
	}

	resp, err := gs.makeAuthenticatedRequest("POST", "credit/wallet/", reqData, nil)
	if err != nil {
		log.Println("Error getting credit wallet response: ", err)
		return nil, fmt.Errorf("failed to get wallet response: %w", err)
	}
	defer resp.Body.Close()

	// Read the raw body
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response body: ", err)
		return nil, err
	}

	// Log the raw response
	log.Println("Raw credit wallet response:", string(rawBody))

	// Decode into your struct
	var result PaymentResponse
	if err := json.Unmarshal(rawBody, &result); err != nil {
		log.Println("Failed to decode credit wallet response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Decoded credit wallet response:", result)
	return &result, nil
}


// VerifyPayment verifies a payment
func (gs *GamingService) VerifyPayment(paymentID string) (map[string]interface{}, error) {
	params := map[string]string{
		"payment_id": paymentID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "verify/payment/", nil, params)
	if err != nil {
		log.Println("Error verifying payment: ", err)
		return nil, fmt.Errorf("failed to verify payment: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode verify payment response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Payment verified successfully: ", result)
	return result, nil
}

// ListPayouts lists all payouts
func (gs *GamingService) ListPayouts() (map[string]interface{}, error) {
	resp, err := gs.makeAuthenticatedRequest("GET", "payouts/", nil, nil)
	if err != nil {
		log.Println("Error listing payouts: ", err)
		return nil, fmt.Errorf("failed to list payouts: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode payouts response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Payouts listed successfully: ", result)
	return result, nil
}

// ListUserPayout lists payouts for a specific user
func (gs *GamingService) ListUserPayout(username string) (map[string]interface{}, error) {
	params := map[string]string{
		"username": username,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "payouts/", nil, params)
	if err != nil {
		log.Println("Error listing user payouts: ", err)
		return nil, fmt.Errorf("failed to list user payouts: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode user payouts response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("User payouts listed successfully: ", result)
	return result, nil
}

// PayoutByAdmin processes payout by admin
func (gs *GamingService) PayoutByAdmin(payoutID string) (map[string]interface{}, error) {
	params := map[string]string{
		"payout_id": payoutID,
	}

	resp, err := gs.makeAuthenticatedRequest("GET", "payouts/", nil, params)
	if err != nil {
		log.Println("Error processing payout: ", err)
		return nil, fmt.Errorf("failed to process payout: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Failed to decode payout processing response: ", err)
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	log.Println("Payout processed successfully: ", result)
	return result, nil
}