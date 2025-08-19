package results

import (
	"net/http"
	"github.com/dblaq/buzzycash/internal/models"
	"github.com/dblaq/buzzycash/internal/utils"
	"github.com/dblaq/buzzycash/pkg/externals"
	"github.com/gin-gonic/gin"
)



func GetWinnerLogsHandler(ctx *gin.Context) {
	gs := externals.NewGamingService()

	logsResponse, err := gs.GetWinnerLogs()
	if err != nil {
		utils.Error(ctx,http.StatusInternalServerError,"Failed to fetch winner logs")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"logsResponse": logsResponse,
		"message":      "Winner logs retrieved successfully",
	})
}

func GetLeaderBoardHandler(ctx *gin.Context) {
	gs := externals.NewGamingService()

	leaderboardResponse, err := gs.GetLeaderBoard()
	if err != nil {
		utils.Error(ctx,http.StatusInternalServerError,"Failed to fetch leaderboard")
		return
	}

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
		utils.Error(ctx,http.StatusInternalServerError,"Failed to fetch user results")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"resultsResponse": resultsResponse,
		"message":         "User results retrieved successfully",
	})
}
