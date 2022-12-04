package middlewares

import (
	"strings"

	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app/otentifikasi"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.JSON(401, gin.H{"error": "Token not found"})
			ctx.Abort()
			return
		}

		err := otentifikasi.ValidateToken(strings.Split(tokenStr, "Bearer ")[1])
		if err != nil {
			ctx.JSON(401, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
