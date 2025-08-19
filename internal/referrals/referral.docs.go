package referral


// @Summary Get referral details
// @Description Fetch referral details for the authenticated user
// @Tags referrals
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Referral details data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Referral details not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /referrals/referral-details [get]
// @Security BearerAuth
func _() {}
