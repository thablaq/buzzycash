package migrations

import (
	"gorm.io/gorm"
	"log"
	// "github.com/dblaq/buzzycash/internal/models"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		// &models.User{},
		// &models.ReferralWallet{},
		// &models.ReferralEarning{},
		// &models.GameHistory{},
		// &models.Notification{},
		// &models.RefreshToken{},
		// &models.Transaction{},
		// &models.UserOtpSecurity{},
		// &models.Role{},
		// &models.RefreshToken{},
		// &models.Admin{},
		// &models.BlacklistedToken{},
	)

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("✅ Database migration completed.")
}