package results

import (
	"net/http"
	"log"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/gin-gonic/gin"
)




func GetWinnerLogsHandler(ctx *gin.Context) {
	gs := externals.NewGamingService()

	logsResponse, err := gs.GetWinnerLogs()
	if err != nil {
		log.Println("Error fetching winner logs:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch winner logs")
		return
	}

	log.Println("Winner logs retrieved successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"logsResponse": logsResponse,
		"message":      "Winner logs retrieved successfully",
	})
}

func GetLeaderBoardHandler(ctx *gin.Context) {
	gs := externals.NewGamingService()

	leaderboardResponse, err := gs.GetLeaderBoard()
	if err != nil {
		log.Println("Error fetching leaderboard:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch leaderboard")
		return
	}

	log.Println("Leaderboard retrieved successfully")
	ctx.JSON(http.StatusOK, gin.H{
		"leaderboardResponse": leaderboardResponse,
		"message":             "Leaderboard retrieved successfully",
	})
}

func GetUserResultsHandler(ctx *gin.Context) {
	gs := externals.NewGamingService()

	currentUser := ctx.MustGet("currentUser").(models.User)
	username := currentUser.PhoneNumber

	resultsResponse, err := gs.GetUserResults(username)
	if err != nil {
		log.Println("Error fetching user results for user:", username, "Error:", err)
		utils.Error(ctx, http.StatusInternalServerError, "Failed to fetch user results")
		return
	}

	log.Println("User results retrieved successfully for user:", username)
	ctx.JSON(http.StatusOK, gin.H{
		"resultsResponse": resultsResponse,
		"message":         "User results retrieved successfully",
	})
}

