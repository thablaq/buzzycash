package results



// @Summary Get winners
// @Description Fetch the list of winners
// @Tags results
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of winners"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /result/winners [get]
// @Security BearerAuth
func _() {}

// @Summary Get leaderboard
// @Description Fetch leaderboard data
// @Tags results
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Leaderboard data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /result/leaderboard [get]
// @Security BearerAuth
func _() {}

// @Summary Get user results
// @Description Fetch results for the authenticated user
// @Tags results
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User results data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User results not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /result/user-results [get]
// @Security BearerAuth
func _() {}
