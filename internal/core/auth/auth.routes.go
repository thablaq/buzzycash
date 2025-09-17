package auth

import (
	"github.com/dblaq/buzzycash/internal/middlewares"
		"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(rg *gin.RouterGroup,db *gorm.DB) {
	authHandler := NewAuthHandler(db)
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.SignUpHandler)
		authRoutes.POST("/login", authHandler.LoginHandler)
		authRoutes.POST("/verify-account", authHandler.VerifyAccountHandler)
		authRoutes.POST("/resend-otp", authHandler.ResendOtpHandler)
		authRoutes.PATCH("/change-password", middlewares.AuthMiddleware, authHandler.ChangePasswordHandler)
		authRoutes.POST("/forgot-password", authHandler.ForgotPasswordHandler)
		authRoutes.POST("/verify-reset-password-otp", authHandler.VerifyPasswordForgotOtpHandler)
		authRoutes.PUT("/reset-password", authHandler.ResetPasswordHandler)
		authRoutes.POST("/logout", authHandler.LogoutHandler)
		authRoutes.POST("/refresh-token", authHandler.RefreshTokenHandler)
	}
}
