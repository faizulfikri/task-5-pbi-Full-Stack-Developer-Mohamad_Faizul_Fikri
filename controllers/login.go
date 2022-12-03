package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app/otentifikasi"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/helper/helperhash"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	user_mod := model.User{}

	json.Unmarshal(body, &user_mod)
	user_mod.Init()
	err := user_mod.Validate("login")

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	user_login := app.UserLogin{}

	err = db.Debug().Table("users").Select("*").Joins("LEFT JOIN photos ON photos.user_id = users.id").
		Where("users.email = ?", user_mod.Email).Find(&user_login).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with email " + user_mod.Email + " not found",
			"data":    nil,
		})
		return
	}

	//verifikasi pass
	err = helperhash.ComparePassword(user_login.Password, user_mod.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": errors.New("Incorect Pass"),
			"data":    nil,
		})
		return
	}

	token, err := otentifikasi.GenerateJWT(user_login.Email, user_login.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserData{
		ID: user_login.ID, Username: user_login.Username, Email: user_login.Email, Token: token,
		Photos: app.Photo{Title: user_login.Title, Caption: user_login.Caption, PhotoUrl: user_login.PhotoUrl},
	}

	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "Success",
		"message": "Login successfully",
		"data":    data,
	})

}

func UpdateUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	user := model.User{}
	err := db.Debug().Where("id = ?", ctx.Param("userId")).First(&user).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with id " + ctx.Param("userId") + " not found",
			"data":    nil,
		})
		return
	}

	body, _ := ioutil.ReadAll(ctx.Request.Body)
	user_mod := model.User{}

	user_mod.ID = user.ID
	json.Unmarshal(body, &user_mod)

	err = user_mod.Validate("login")
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user_mod.HashPassword()

	//update user
	err = db.Debug().Model(&user).Save(&user_mod).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	data := app.UserRegister{
		ID:       user_mod.ID,
		Username: user_mod.Username,
		Email:    user_mod.Email,
		CreateAt: user_mod.CreateAt,
		UpdateAt: user_mod.UpdateAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Error",
		"message": "User updated succesfully",
		"data":    data,
	})
}

func DeleteUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	user_mod := model.User{}
	err := db.Debug().Where("id = ?", ctx.Param("userId")).First(&user_mod).Error

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": "User with id " + ctx.Param("userId") + " not found",
			"data":    nil,
		})
		return
	}

	//Delete user
	err = db.Debug().Delete(&user_mod).Error
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Response success
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "User deleted succesfully",
		"data":    nil,
	})
}
