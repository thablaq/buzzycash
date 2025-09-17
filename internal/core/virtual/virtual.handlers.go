package virtual

import (
	"net/http"
     "log"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/external/gaming"
	"github.com/gin-gonic/gin"
)



func GetVirtualGamesHandler(ctx *gin.Context) {
	log.Println("Fetching virtual games...")

	gs := gaming.GMInstance()
	gamesResponse, err := gs.GetVirtualGames()
	if err != nil {
		log.Println("Failed to fetch virtual games:", err)
		utils.Error(ctx,http.StatusInternalServerError, "Failed to fetch virtual games")
		return
	}

	log.Println("Virtual games retrieved successfully.")
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gamesResponse,
		"message": "Virtual games retrieved successfully",
	})
}


func StartVirtualGameHandler(ctx *gin.Context) {
	var req StartGameRequest

	
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, utils.ValidationErrorToJSON(err))
		return
	}

	currentUser := ctx.MustGet("currentUser").(models.User)
	username := currentUser.PhoneNumber
	log.Printf("Validated request data and extracted username: %s\n", username)

	gs := gaming.GMInstance()
	gameData, err := gs.StartVirtualGame(req.GameType, username)
	if err != nil {
		log.Println("Failed to start virtual game:", err)
		utils.Error(ctx,http.StatusInternalServerError, "Failed to start virtual game")
		return
	}

	log.Println("Virtual game started successfully.")
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":   gameData,
		"message": "Virtual game started successfully",
	})
}
