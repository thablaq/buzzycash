package tickets

import (
	"fmt"
	"net/http"
	"time"

	"log"

	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/helpers"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/gaming"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BuyGameTicketHandler(ctx *gin.Context) {
	var req BuyTicketRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	log.Println("Attempting to validate request struct")
	if err := req.Validate(); err != nil {
    utils.Error(ctx, http.StatusBadRequest, err.Error())
    return
}

	currentUser := ctx.MustGet("currentUser").(models.User)
	userID := currentUser.ID
	username := currentUser.PhoneNumber

	log.Println("Attempting to purchase ticket for user: %s, game_id: %s, quantity: %d, amount: %.2f",
		username, req.GameID, req.Quantity, req.AmountPaid)

	transactionTxRef := helpers.GenerateTransactionReference()
	log.Println("[transactionTxRef] ✅ Unique transaction ref generated: %s", transactionTxRef)

	gs := gaming.GMInstance()
	buyResponse, err := gs.BuyTicket(req.GameID, username, req.Quantity, req.AmountPaid)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to buy ticket")
		return
	}

	log.Println("Received response for ticket purchase: %+v", buyResponse)

	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// Insert tickets
		for _, ticketID := range buyResponse.TicketIDs {
			ticket := models.TicketPurchase{
				ID:          ticketID,
				UserID:      userID,
				TotalAmount: req.AmountPaid,
				UnitPrice:   req.AmountPaid / int64(req.Quantity),
				Quantity:    1,
				Currency:    "NGN",
				PurchasedAt: time.Now(),
			}
			if err := tx.Create(&ticket).Error; err != nil {
				return err
			}
		}

		// Create transaction history linked to first ticket
		history := models.TransactionHistory{
			Amount:           req.AmountPaid,
			UserID:               userID,
			PaymentStatus:        models.Successful,
			PaymentMethod:        models.Wallet,
			TransactionReference: transactionTxRef,
			TransactionType:      models.Debit,
			Category:             models.Ticket,
			Currency:             "NGN",
			TicketPurchaseID:     buyResponse.TicketIDs[0],
			Metadata:             map[string]interface{}{"ticketIds": buyResponse.TicketIDs, "gameId": req.GameID},
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to save ticket purchase")
		return
	}

	
	// notification := models.Notification{
	// 	UserID:  userID,
	// 	Title:   "Ticket Purchase Successful",
	// 	Message: fmt.Sprintf("Your request to purchase ticket with ₦%.2f has been successful.", req.AmountPaid),
	// 	Type:    models.Ticket,
	// 	IsRead:  false,
	// }
	// config.DB.Create(&notification)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"tickets": buyResponse,
		"message": "Ticket purchased successfully",
	})
}

func GetUserGameTicketsHandler(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	username := currentUser.PhoneNumber

	log.Printf("Fetching tickets for user: %s", username)

	gs := gaming.GMInstance()
	ticketsResult, err := gs.GetUserTickets(username)
	if err != nil {
		log.Printf("Error fetching tickets for user %s: %v", username, err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch user tickets")
		return
	}

	log.Printf("Successfully fetched tickets for user %s: %+v", username, ticketsResult)

	// Extract fields from the map
	gameID := ticketsResult["game_id"]
	purchasedAt := ticketsResult["purchased_at"]
	status := ticketsResult["status"]
	ticketsResponse := ticketsResult["tickets"]

	ctx.JSON(http.StatusOK, gin.H{
			"game_id":      gameID,
			"purchased_at": purchasedAt,
			"status":       status,
			"tickets":      ticketsResponse,
		"message": "User games retrieved successfully",
	})
}

func GetAllGamesHandler(ctx *gin.Context) {
	log.Println("Fetching all games")

	gs := gaming.GMInstance()
	gameResults, err := gs.GetGames()
	if err != nil {
		log.Printf("Error retrieving games: %v", err)
		utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve games: %v", err))
		return
	}

	log.Printf("Successfully retrieved games: %+v", gameResults)

	// Return the successful response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Games retrieved successfully.",
		"results":   gameResults,
	})
}

func CreateGameHandler(ctx *gin.Context) {
	var req CreateGameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request payload: %v", err)
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	// if err := utils.Validate.Struct(req); err != nil {
	// 	log.Printf("Validation error: %v", err)
	// 	utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
	// 	return
	// }

	if req.WinningPercentage < 0 || req.WinningPercentage > 100 {
		log.Printf("Invalid winning_percentage: %d", req.WinningPercentage)
		utils.Error(ctx, http.StatusBadRequest, "winning_percentage must be between 0 and 100")
		return
	}

	log.Printf("Creating game with request: %+v", req)

	gs := gaming.GMInstance()
	gameResponse, err := gs.CreateGames(
		req.GameName,
		req.Amount,
		req.DrawInterval,
		req.WinningPercentage,
		req.MaxWinners,
		req.Date,
		req.WeightedDistribution,
	)
	if err != nil {
		log.Printf("Error creating game: %v", err)
		utils.Error(ctx, http.StatusInternalServerError, fmt.Sprintf("Failed to create game: %v", err))
		return
	}

	log.Printf("Game created successfully: %+v", gameResponse)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Game created successfully",
		"data":    gameResponse,
	})
}
