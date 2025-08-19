package wallets



// @Summary Credit wallet
// @Description Generate a payment link to credit the user's wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body creditWalletRequest true "Wallet credit request"
// @Success 201 {object} map[string]interface{} "Payment link generated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to generate payment link"
// @Router /wallet/request-link [post]
// @Security BearerAuth
func _() {}


// @Summary Get wallet balance
// @Description Retrieve the authenticated user's wallet balance
// @Tags wallet
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User wallet balance"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to fetch wallet balance"
// @Router /wallet/get-wallet [get]
// @Security BearerAuth
func _() {}
