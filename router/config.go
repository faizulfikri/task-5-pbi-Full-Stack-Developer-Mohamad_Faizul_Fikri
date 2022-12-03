package router

import (
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		ctx.Set("db", db)
	})

	//User
	r.POST("/users/register", controllers.RegisterNewUser)
	r.POST("/users/login", controllers.Login)
	r.GET("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUser)

	r.GET("/photos/")
	return r
}
