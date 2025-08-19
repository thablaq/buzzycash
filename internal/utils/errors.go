package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"
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

// Error helper
func Error(ctx *gin.Context, statusCode int, message interface{}) {
	err := &AppError{
		StatusCode: statusCode,
		Message:    message,
	}
	err.Write(ctx)
}



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
		errors["error"] = err.Error()
	}
	return errors
}



