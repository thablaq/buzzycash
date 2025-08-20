package tickets

// @Summary Purchase game ticket
// @Description Buy a ticket for a specific game
// @Tags tickets
// @Accept json
// @Produce json
// @Param request body BuyTicketRequest true "Ticket purchase details"
// @Success 201 {object} map[string]interface{} "Ticket purchased successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /ticket/purchase-ticket [post]
// @Security BearerAuth
func _() {}

// @Summary Get user tickets
// @Description Retrieve all purchased tickets for the authenticated user
// @Tags tickets
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of user tickets"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "No tickets found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /ticket/get-tickets [get]
// @Security BearerAuth
func _() {}


// @Summary Get all games
// @Description Retrieve a list of all available games from the gaming service
// @Tags gaming
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all games"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /gaming [get]
// @Security BearerAuth
func _() {}
