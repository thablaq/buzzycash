package virtual

import (
	"net/http"
     "fmt"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/gin-gonic/gin"
)
// Handler to fetch virtual games
func GetVirtualGamesHandler(ctx *gin.Context) {
	fmt.Println("Fetching virtual games...")

	gs := externals.NewGamingService()
	gamesResponse, err := gs.GetVirtualGames()
	if err != nil {
		fmt.Println("Failed to fetch virtual games:", err)
		utils.Error(ctx,http.StatusInternalServerError, "Failed to fetch virtual games")
		return
	}

	fmt.Println("Virtual games retrieved successfully.")
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"gamesResponse": gamesResponse},
		"message": "Virtual games retrieved successfully",
	})
}

// Handler to start a virtual game
func StartVirtualGameHandler(ctx *gin.Context) {
	var req StartGameRequest

	// Bind and validate request body
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Validate(); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Extract current user's phone number from context
	currentUser := ctx.MustGet("currentUser").(models.User)
	username := currentUser.PhoneNumber
	fmt.Printf("Validated request data and extracted username: %s\n", username)

	gs := externals.NewGamingService()
	gameData, err := gs.StartVirtualGame(req.GameType, username)
	if err != nil {
		fmt.Println("Failed to start virtual game:", err)
		utils.Error(ctx,http.StatusInternalServerError, "Failed to start virtual game")
		return
	}

	fmt.Println("Virtual game started successfully.")
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    gin.H{"gameData": gameData},
		"message": "Virtual game started successfully",
	})
}
