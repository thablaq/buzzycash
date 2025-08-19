package virtual


// @Summary Start a virtual game
// @Description Starts a new virtual game session for the authenticated user
// @Tags virtual
// @Accept json
// @Produce json
// @Param request body StartGameRequest true "Game start data"
// @Success 201 {object} map[string]interface{} "Game started successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request payload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to start game"
// @Router /virtual/start-game [post]
// @Security BearerAuth
func _() {}


// @Summary Get virtual games
// @Description Retrieve a list of available virtual games
// @Tags virtual
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of virtual games"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to fetch games"
// @Router /virtual/get-games [get]
// @Security BearerAuth
func _() {}
