package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/dblaq/buzzycash/internal/middlewares"
)

func AuthRoutes(rg *gin.RouterGroup) {
	authRoutes := rg.Group("/auth")
	{
		authRoutes.POST("/register", SignUpHandler)
		authRoutes.POST("/login", LoginHandler)
		authRoutes.POST("/verify-account", VerifyAccountHandler)
		authRoutes.POST("/resend-otp", ResendOtpHandler)
		authRoutes.PATCH("/change-password", middlewares.AuthMiddleware, ChangePasswordHandler)
		authRoutes.POST("/forgot-password", ForgotPasswordHandler)
		authRoutes.POST("/verify-reset-password-otp", VerifyPasswordForgotOtpHandler)
		authRoutes.PUT("/reset-password", ResetPasswordHandler)
		authRoutes.POST("/logout", LogoutHandler)
		authRoutes.POST("/refresh-token", RefreshTokenHandler)
	}
}
