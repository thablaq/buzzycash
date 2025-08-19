package utils

import (
	"github.com/dblaq/buzzycash/internal/config"
	"github.com/dblaq/buzzycash/internal/models"
	"time"
    "fmt"
    "gorm.io/gorm/clause"
    "errors"
    "gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

var (
	AccessTokenTTL  = time.Hour * 24 * 3
	RefreshTokenTTL = time.Hour * 24 * 7
)


func GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JwtAccessSecret))
}


func GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(RefreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JwtAccessSecret))
}


func VerifyJWTRefreshToken(tokenStr string) (string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(config.AppConfig.JwtRefreshSecret), nil
    })

    if err != nil || !token.Valid {
        return "", fmt.Errorf("invalid or expired refresh token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", fmt.Errorf("invalid token claims")
    }

    // Get user ID safely
    userID, ok := claims["user_id"].(string)
    if !ok || userID == "" {
        return "", fmt.Errorf("user_id not found in token")
    }

    // // Validate UUID format if needed
    // if _, err := uuid.Parse(userID); err != nil {
    //     return "", fmt.Errorf("invalid user ID format")
    // }

    return userID, nil
}


func DecodeToken(tokenStr string) (map[string]interface{}, error) {
    token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
    if err != nil {
        return nil, fmt.Errorf("invalid token: %v", err)
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }
    
    return claims, nil
}



// BlacklistToken adds the given token to the blacklist with an expiration time
func BlacklistToken(token string, expireAt time.Time) error {
    blacklisted := models.BlacklistedToken{
        Token:    token,
        ExpiresAt: expireAt,
    }

    // Use Create with OnConflict to prevent duplicates
    return config.DB.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "token"}},
        DoNothing: true,
    }).Create(&blacklisted).Error
}

// IsTokenBlacklisted checks if a token is blacklisted
func IsTokenBlacklisted(token string) (bool, error) {
    var b models.BlacklistedToken
    err := config.DB.Where("token = ?", token).First(&b).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil
        }
        return false, err
    }

    // Remove expired tokens automatically
    if time.Now().After(b.ExpiresAt) {
        if delErr := config.DB.Delete(&b).Error; delErr != nil {
            return false, delErr
        }
        return false, nil
    }

    return true, nil
}

