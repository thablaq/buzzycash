// package utils

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/go-playground/validator/v10"
// 	"strings"
// )


// type AppError struct {
// 	StatusCode int         `json:"-"`
// 	Message    interface{} `json:"message"` 
// }

// func (e *AppError) ToJSON() map[string]interface{} {
// 	return map[string]interface{}{
// 		"message": e.Message,
// 	}
// }

// func (e *AppError) Write(ctx *gin.Context) {
// 	ctx.JSON(e.StatusCode, e.ToJSON())
// }

// // Error helper
// func Error(ctx *gin.Context, statusCode int, message interface{}) {
// 	err := &AppError{
// 		StatusCode: statusCode,
// 		Message:    message,
// 	}
// 	err.Write(ctx)
// }



// func ValidationErrorToJSON(err error) map[string]string {
// 	errors := make(map[string]string)
// 	if errs, ok := err.(validator.ValidationErrors); ok {
// 		for _, e := range errs {
// 			field := strings.ToLower(e.Field())
// 			var msg string

// 			switch e.Tag() {
// 			case "required":
// 				msg = field + " is required"
// 			case "min":
// 				msg = field + " must be at least " + e.Param() + " characters long"
// 			case "max":
// 				msg = field + " must be at most " + e.Param() + " characters long"
// 			case "eqfield":
// 				msg = field + " must match " + strings.ToLower(e.Param())
// 			case "email":
// 				msg = "invalid email format"
// 			default:
// 				msg = "invalid value for " + field
// 			}

// 			errors[field] = msg
// 		}
// 	} else {
// 		errors["error"] = err.Error()
// 	}
// 	return errors
// }




package utils

import (
	"log"
	// "regexp"
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

// Database error patterns that should be sanitized
var dbErrorPatterns = []string{
	"duplicate key value",
	"violates unique constraint",
	"violates foreign key constraint", 
	"violates check constraint",
	"violates not-null constraint",
	"column cannot be null",
	"table doesn't exist",
	"unknown column",
	"syntax error",
	"connection refused",
	"gorm",
	"sql:",
	"database",
	"constraint",
	"relation",
	"does not exist",
	"already exists",
	"deadlock",
	"timeout",
}

// Enhanced Error helper that automatically sanitizes database errors in production
func Error(ctx *gin.Context, statusCode int, message interface{}) {
	// Log the original message for debugging
	log.Printf("Error occurred [%d]: %v", statusCode, message)
	
	// If we're in production and status code indicates server error, sanitize the message
	if config.AppConfig.Env == "production" && statusCode >= 500 {
		message = sanitizeErrorMessage(message)
	}
	
	// For 4xx errors in production, also check for database leaks
	if config.AppConfig.Env == "production" && statusCode >= 400 && statusCode < 500 {
		message = sanitizeDatabaseLeaks(message)
	}
	
	err := &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
	err.Write(ctx)
}

func sanitizeErrorMessage(message interface{}) interface{} {
	var messageStr string
	
	switch v := message.(type) {
	case string:
		messageStr = v
	case error:
		messageStr = v.Error()
	case map[string]string:
		// Handle validation errors - sanitize each field
		sanitizedMap := make(map[string]string)
		for key, value := range v {
			sanitizedMap[key] = sanitizeString(value)
		}
		return sanitizedMap
	case map[string]interface{}:
		// Handle generic maps
		sanitizedMap := make(map[string]interface{})
		for key, value := range v {
			if strValue, ok := value.(string); ok {
				sanitizedMap[key] = sanitizeString(strValue)
			} else {
				sanitizedMap[key] = value
			}
		}
		return sanitizedMap
	default:
		return "Internal server error"
	}
	
	return sanitizeString(messageStr)
}

func sanitizeDatabaseLeaks(message interface{}) interface{} {
	messageStr, ok := message.(string)
	if !ok {
		return message
	}
	
	// Check if message contains database-related terms
	lowerMsg := strings.ToLower(messageStr)
	for _, pattern := range dbErrorPatterns {
		if strings.Contains(lowerMsg, pattern) {
			log.Printf("Detected potential database leak in 4xx error: %s", messageStr)
			return "Bad request"
		}
	}
	
	return message
}

func sanitizeString(input string) string {
	if input == "" {
		return input
	}
	
	lowerInput := strings.ToLower(input)
	
	// Check for database error patterns
	for _, pattern := range dbErrorPatterns {
		if strings.Contains(lowerInput, pattern) {
			// Return appropriate generic message based on the type of error
			if strings.Contains(lowerInput, "duplicate") || strings.Contains(lowerInput, "unique") {
				return "Resource already exists"
			}
			if strings.Contains(lowerInput, "foreign key") || strings.Contains(lowerInput, "constraint") {
				return "Data validation failed"
			}
			if strings.Contains(lowerInput, "not-null") || strings.Contains(lowerInput, "cannot be null") {
				return "Required field missing"
			}
			if strings.Contains(lowerInput, "connection") || strings.Contains(lowerInput, "timeout") {
				return "Service temporarily unavailable"
			}
			if strings.Contains(lowerInput, "deadlock") {
				return "Please try again"
			}
			
			// Default fallback
			return "Internal server error"
		}
	}
	
	// Also check for SQL injection-like patterns or sensitive info
	sensitivePatterns := []string{
		"select ", "insert ", "update ", "delete ", "drop ", "alter ",
		"create table", "database", "sql", "mysql", "postgresql", "gorm",
	}
	
	for _, pattern := range sensitivePatterns {
		if strings.Contains(lowerInput, pattern) {
			return "Internal server error"
		}
	}
	
	return input
}

// Keep your existing ValidationErrorToJSON function as is
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
		// Sanitize the error message if we're in production
		errorMsg := err.Error()
		if config.AppConfig.Env == "production" {
			errorMsg = sanitizeString(errorMsg)
		}
		errors["error"] = errorMsg
	}
	return errors
}