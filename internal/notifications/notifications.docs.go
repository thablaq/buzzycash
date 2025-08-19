package notifications


// @Summary Get all notifications
// @Description Retrieve all notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {array} map[string]interface{} "List of notifications"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to fetch notifications"
// @Router /notification/ [get]
// @Security BearerAuth
func _() {}


// @Summary Get unread notifications count
// @Description Retrieve the count of unread notifications for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Unread notifications count"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to fetch unread count"
// @Router /notification/unread [get]
// @Security BearerAuth
func _() {}


// @Summary Mark a notification as read
// @Description Mark a specific notification as read using its ID
// @Tags notifications
// @Accept json
// @Produce json
// @Param notificationId path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification marked as read"
// @Failure 400 {object} map[string]interface{} "Invalid notification ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Notification not found"
// @Failure 500 {object} map[string]interface{} "Failed to mark notification as read"
// @Router /notification/{notificationId}/read [patch]
// @Security BearerAuth
func _() {}


// @Summary Mark all notifications as read
// @Description Mark all notifications as read for the authenticated user
// @Tags notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "All notifications marked as read"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Failed to mark all as read"
// @Router /notification/read-all [patch]
// @Security BearerAuth
func _() {}
