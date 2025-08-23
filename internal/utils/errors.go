package utils

import (
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/dblaq/buzzycash/internal/config"
)

type AppError struct {
	StatusCode int         `json:"-"`
	Message    interface{} `json:"message"`
}

func (e *AppError) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"message": e.Message,
	}
}

func (e *AppError) Write(ctx *gin.Context) {
	ctx.JSON(e.StatusCode, e.ToJSON())
}


var pgErrors = map[string]string{
	"22p02": "Invalid data format provided",
	"23505": "Resource already exists",
	"23503": "Invalid reference provided",
	"23514": "Data validation failed",
	"23502": "Required field missing",
	"42p01": "Internal server error",
	"42703": "Internal server error",
}

// Error helper with automatic database error sanitization
func Error(ctx *gin.Context, statusCode int, message interface{}) {
	log.Printf("Error [%d]: %v", statusCode, message)
	
	// Sanitize database errors in production for 4xx/5xx status codes
	isProduction := config.AppConfig.Env == "production"
	if isProduction && statusCode >= 400 {
		message = sanitize(message)
	}
	
	(&AppError{statusCode, message}).Write(ctx)
}

// Sanitize database errors
func sanitize(message interface{}) interface{} {
	switch v := message.(type) {
	case string:
		return sanitizeString(v)
	case error:
		return sanitizeString(v.Error())
	case map[string]string:
		result := make(map[string]string)
		for k, val := range v {
			result[k] = sanitizeString(val)
		}
		return result
	default:
		return "Internal server error"
	}
}

func sanitizeString(input string) string {
	if input == "" {
		return input
	}
	
	lower := strings.ToLower(input)
	
	// Check PostgreSQL error codes
	for code, msg := range pgErrors {
		if strings.Contains(lower, code) {
			return msg
		}
	}
	
	// Check SQLSTATE pattern
	if strings.Contains(lower, "sqlstate") {
		if matches := regexp.MustCompile(`sqlstate\s+(\w+)`).FindStringSubmatch(lower); len(matches) > 1 {
			if msg, exists := pgErrors[matches[1]]; exists {
				return msg
			}
		}
		return "Internal server error"
	}
	
	return input
}

// Validation errors to JSON
func ValidationErrorToJSON(err error) map[string]string {
	errors := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			var msg string
			switch e.Tag() {
			case "required":
				msg = field + " is required"
			case "min":
				msg = field + " must be at least " + e.Param() + " characters long"
			case "max":
				msg = field + " must be at most " + e.Param() + " characters long"
			case "eqfield":
				msg = field + " must match " + strings.ToLower(e.Param())
			case "email":
				msg = "invalid email format"
			default:
				msg = "invalid value for " + field
			}
			errors[field] = msg
		}
	} else {
		errorMsg := err.Error()
		if config.AppConfig.Env == "production" {
			errorMsg = sanitizeString(errorMsg)
		}
		errors["error"] = errorMsg
	}
	return errors
}