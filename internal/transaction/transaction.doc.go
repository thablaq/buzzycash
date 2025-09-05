package transaction

// @Summary Get all transactions
// @Description Retrieve a list of all transactions, optionally filtered by parameters like status or date
// @Tags transactions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all transactions"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/history [get]
// @Security BearerAuth
func _() {}

// @Summary Get transaction by ID
// @Description Retrieve the details of a single transaction using its ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]interface{} "Transaction details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Transaction not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/{id} [get]
// @Security BearerAuth
func _() {}

// @Summary Search transactions
// @Description Search through transactions using keywords like reference, amount, or status
// @Tags transactions
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "No transactions found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/search [get]
// @Security BearerAuth
func _() {}
