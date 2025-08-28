package notifications

import (
	"fmt"
	// "time"
	"github.com/dblaq/buzzycash/internal/models"
)


func mapNotificationToResponse(n models.Notification) NotificationResponse {
	return NotificationResponse{
		Title:       n.Title,
		Subtitle:    n.Subtitle,
		Amount:      fmt.Sprintf("%.2f %s", n.Amount, n.Currency),
		Status:      n.Status,
		CreatedAt:   n.CreatedAt,
		DisplayTime: n.CreatedAt.Format("Jan 02 2006 @03:04pm"),
	}
}





// You can customize titles/subtitles per category/method
func BuildTxNotifContent(h models.TransactionHistory) (title, subtitle string) {
	switch h.Category {
	case models.Cashout, models.WithdrawRequest:
		return "Cashout Successful", "Your withdrawal via bank was successful"
	case models.Deposit:
		return "Deposit Successful", "You have successfully deposited into your wallet."
	case models.PrizeMoney:
		return "Wallet Credited", "Your wallet has been credited for prize money."
	case models.Ticket:
		return "Ticket Purchased", "You purchased ticket"
	default:
		// fallback by payment status or method
		if h.PaymentStatus == models.Successful {
			return "Payment Successful", "Your transaction was successful."
		}
		return "Transaction Update", "Your transaction status changed."
	}
}