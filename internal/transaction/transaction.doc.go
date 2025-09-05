package transaction

// @Summary Get all transactions
// @Description Retrieve a list of all transactions, optionally filtered by parameters like status or date
// @Tags transactions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param payment_status query string false "Filter by payment status"
// @Param payment_type query string false "Filter by payment type"
// @Param category query string false "Filter by category"
// @Param transaction_type query string false "Filter by transaction type"
// @Param payment_method query string false "Filter by payment method"
// @Param currency query string false "Filter by currency"
// @Success 200 {object} TransactionHistoryResponseList "List of all transactions"
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
// @Description Search through transactions using keywords like reference, email, category, or status
// @Tags transactions
// @Accept json
// @Produce json
// @Param search query string true "Search query"
// @Param page query int false "Page number (default: 1)"
// @Success 200 {object} TransactionHistoryResponseList "Search results"
// @Failure 400 {object} map[string]interface{} "Bad request (e.g., empty search query)"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "No transactions found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /transactions/search [get]
// @Security BearerAuth
func _() {}


