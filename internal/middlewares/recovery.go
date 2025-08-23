package middlewares

import (
	"log"
	"github.com/dblaq/buzzycash/internal/config"
	"net/http"
	"github.com/gin-gonic/gin"
)

func RecoveryAndErrorMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        defer func() {
            if r := recover(); r != nil {
                // Panic caught
                log.Printf("panic recovered: %v", r)

                // Respond safely
                ctx.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
                ctx.Abort()
            }
        }()

        ctx.Next()

        // Check for collected errors (e.g., DB errors)
        if len(ctx.Errors) > 0 {
            err := ctx.Errors[0].Err
            if config.AppConfig.Env=="production" {
                log.Println("Error occurred:", err)
                ctx.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
            } else {
                ctx.JSON(http.StatusInternalServerError, gin.H{
                    "error": err.Error(),
                })
            }
            ctx.Abort()
        }
    }
}
