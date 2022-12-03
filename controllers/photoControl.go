package controllers

import (
	"net/http"

	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPhoto(ctx *gin.Context) {
	photos := []model.Photo{}

	db := ctx.MustGet("db").(*gorm.DB)
	if err := db.Debug().Model(&model.Photo{}).Limit(100).Find(&photos).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	if len(photos) > 0 {
		for p := range photos {
			user := model.User{}
			err := db.Model(&model.User{}).Where("id = ?", photos[p].UserID).Take(&user).Error

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": err.Error(),
					"data":    nil,
				})
				return
			}

			photos[p].Owner = app.Owner{
				ID: user.ID, Username: user.Username, Email: user.Email,
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Data retrieved successfully",
		"data":    photos,
	})
}
