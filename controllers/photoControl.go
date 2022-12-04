package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app/otentifikasi"
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
			err := db.Model(&model.User{}).Where("id = ?", photos[p].ID).Take(&user).Error

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

func CreatePhoto(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	tokenStr := ctx.GetHeader("Authorization")

	if tokenStr == "" {
		ctx.JSON(401, gin.H{"error": "T"})
		return
	}

	email, err := otentifikasi.GetMail(strings.Split(tokenStr, "Bearer ")[1])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	user_has_log := model.User{}
	err = db.Debug().Where("email = ?", email).First(&user_has_log).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
			"data":    nil,
		})
		return
	}

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	inp_photo := model.Photo{}
	err = json.Unmarshal(body, &inp_photo)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	inp_photo.Init()
	inp_photo.UserID = user_has_log.ID
	inp_photo.Owner = app.Owner{
		ID:       user_has_log.ID,
		Username: user_has_log.Username,
		Email:    user_has_log.Email,
	}

	err = inp_photo.Validate("upload")
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	var old_photo model.Photo
	err = db.Debug().Model(&model.Photo{}).Where("user_id = ?", user_has_log.ID).Find(&old_photo).Error
	if err != nil {
		if err.Error() == "Data not found" {
			err = db.Debug().Create(&inp_photo).Error //Create photo to database
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Error",
					"message": err,
					"data":    nil,
				})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"status":  "Success",
				"message": "Photo uploaded successfully",
				"data":    inp_photo,
			})
			return
		}
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})

		return
	}

	inp_photo.ID = old_photo.ID
	err = db.Debug().Model(&old_photo).Updates(&inp_photo).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": err,
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo changed successfully",
		"data":    inp_photo,
	})
}

func DeletePhoto(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	tokenStr := ctx.GetHeader("Authorization")

	if tokenStr == "" {
		ctx.JSON(401, gin.H{"error": "T"})
		return
	}

	email, err := otentifikasi.GetMail(strings.Split(tokenStr, "Bearer ")[1])
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	var user_has_log model.User
	if err := db.Debug().Where("email = ?", email).First(&user_has_log).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + email + " not found",
			"data":    nil})
		return
	}

	var photo model.Photo
	if err = db.Debug().Where("id = ?", ctx.Param("photoId")).First(&photo).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "Photo not found",
			"data":    nil,
		})
		return
	}

	if user_has_log.ID != photo.UserID {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "You can't delete photo of another user",
			"data":    nil,
		})
		return
	}

	err = db.Debug().Delete(&photo).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Photo deleted successfully",
		"data":    nil,
	})

}
