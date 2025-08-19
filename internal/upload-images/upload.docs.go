package uploadimages

// @Summary Upload user avatar
// @Description Upload a profile avatar image for the authenticated user
// @Tags uploads
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} map[string]interface{} "File uploaded successfully"
// @Failure 400 {object} map[string]interface{} "Invalid file upload"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /upload/user [post]
// @Security BearerAuth
func _() {}
