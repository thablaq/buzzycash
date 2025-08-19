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
    authHeader := ctx.GetHeader("Authorization")
    if authHeader == "" {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    if tokenString == authHeader {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
        return
    }

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.AppConfig.JwtAccessSecret), nil
    })

    if err != nil || !token.Valid {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
        return
    }

    // Check blacklist
    var blacklisted models.BlacklistedToken
    if err := config.DB.First(&blacklisted, "token = ?", tokenString).Error; err == nil {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token blacklisted"})
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
        return
    }

    // Safe expiration check
    if expRaw, ok := claims["exp"].(float64); ok {
        if time.Now().Unix() > int64(expRaw) {
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
            return
        }
    } else {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token missing expiration"})
        return
    }

    
    userID, ok := claims["user_id"].(string)
    if !ok || userID == "" {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
        return
    }

    var user models.User
    if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    ctx.Set("currentUser", user)
    ctx.Next()
}