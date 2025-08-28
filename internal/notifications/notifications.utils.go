package notifications

import (
	"fmt"
	// "time"
	"github.com/dblaq/buzzycash/internal/models"
)


func mapTransactionToNotification(tx models.TransactionHistory) NotificationResponse {
	return NotificationResponse{
		Title:       string(tx.PaymentStatus),                        // e.g. "Cashout Successful"
		Subtitle:    string(tx.PaymentStatus),                    // "Your withdrawal via bank was successful"
		Amount:      fmt.Sprintf("%.2f %s", tx.Amount, tx.Currency),  // 2000.00 NGN
		Status:      string(tx.TransactionType),                      // "Ticket Purchased"
		CreatedAt:   tx.CreatedAt,
		DisplayTime: tx.CreatedAt.Format("Jan 02 2006 @03:04pm"),
	}
}

// func mapGameToNotification(game models.GameHistory) NotificationResponse {
// 	return NotificationResponse{
// 		Title:       game.Name,                                       // e.g. "Daily ChopChop"
// 		Subtitle:    game.Status,                                     // e.g. "Ongoing", "Ended"
// 		Amount:      formatGameAmount(game),                          // "__ NGN" if not set, or "2000.00 NGN"
// 		Status:      game.Result,                                     // "You Won", "You Lost", "Ongoing"
// 		CreatedAt:   game.CreatedAt,
// 		DisplayTime: game.CreatedAt.Format("Jan 02 2006 @03:04pm"),
// 	}
// }

// func formatGameAmount(g models.GameHistory) string {
// 	if g.Amount == 0 {
// 		return "__ NGN"
// 	}
// 	return fmt.Sprintf("%.2f NGN", g.Amount)
// }

// func buildTransactionSubtitle(tx models.TransactionHistory) string {
// 	switch tx.PaymentStatus {
// 	case "Cashout Successful":
// 		return "Your withdrawal via bank was successful"
// 	case "Wallet Credited":
// 		return "Your wallet has been credited for prize money."
// 	case "Deposit Successful":
// 		return "You have successfully deposited into wallet."
// 	default:
// 		return "Transaction processed"
// 	}
// }



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