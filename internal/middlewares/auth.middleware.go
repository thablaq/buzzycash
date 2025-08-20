package middlewares

import (
	"fmt"
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(ctx *gin.Context) {
	// Extract token
	authHeader := ctx.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if authHeader == "" || tokenString == authHeader {
		abortWithError(ctx, "Invalid or missing authorization token")
		return
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.AppConfig.JwtAccessSecret), nil
	})

	if err != nil || !token.Valid {
		abortWithError(ctx, "Invalid or expired token")
		return
	}

	// Check blacklist
	var blacklisted models.BlacklistedToken
	if err := config.DB.First(&blacklisted, "token = ?", tokenString).Error; err == nil {
		abortWithError(ctx, "Token blacklisted")
		return
	}

	// Extract and validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		abortWithError(ctx, "Invalid token claims")
		return
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); !ok || time.Now().Unix() > int64(exp) {
		abortWithError(ctx, "Token expired")
		return
	}

	// Get user
	userID, ok := claims["user_id"].(string)
	if !ok || userID == "" {
		abortWithError(ctx, "Invalid user ID")
		return
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
		abortWithError(ctx, "User not found")
		return
	}

	ctx.Set("currentUser", user)
	ctx.Next()
}

func abortWithError(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": message})
}