package notifications


import "time"


type NotificationResponse struct {
	Title       string    `json:"title"`       // e.g. "Cashout Successful", "Daily ChopChop"
	Subtitle    string    `json:"subtitle"`    // e.g. "Your withdrawal via bank was successful", "Ongoing"
	Amount      string    `json:"amount"`      // formatted with .00 NGN or "__ NGN" for ongoing
	Status      string    `json:"status"`      // e.g. "You Won", "You Lost", "Ongoing"
	CreatedAt   time.Time `json:"created_at"`  // actual date/time
	DisplayTime string    `json:"display_time"`// formatted "Oct 26 2024 @10:58pm"
}
