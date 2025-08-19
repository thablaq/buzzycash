package profile


// @Summary Get user profile
// @Description Fetch the authenticated user’s profile
// @Tags profile
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Profile data"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profile/get-profile [get]
// @Security BearerAuth
func _() {}


// @Summary Create user profile
// @Description Create a new profile for the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Param request body CreateProfileRequest true "Profile creation data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 409 {object} map[string]interface{} "Profile already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profile/create-profile [post]
// @Security BearerAuth
func _() {}



// @Summary Update user profile
// @Description Update fields of the authenticated user’s profile
// @Tags profile
// @Accept json
// @Produce json
// @Param request body ProfileUpdateRequest true "Profile update data"
// @Success 200 {object} map[string]interface{} "Profile updated successfully"
// @Failure 400 {object} map[string]interface{} "Validation error"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Profile not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /profile/update-profile [patch]
// @Security BearerAuth
func _() {}
