package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/app"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// function for handle register request
func RegisterNewUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)

	//read body req
	body, _ := ioutil.ReadAll(ctx.Request.Body)

	//Convert json to obj
	user_mod := model.User{}
	json.Unmarshal(body, &user_mod)

	user_mod.Init()                      //inisialisasi
	err := user_mod.Validate("register") //Validate user

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	//hash pass
	err = user_mod.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Debug().Create(&user_mod).Error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Error",
			"message": err.Error(),
			"data":    nil,
		})
	}

	data := app.UserRegister{
		ID:       user_mod.ID,
		Username: user_mod.Username,
		Email:    user_mod.Email,
		CreateAt: user_mod.CreateAt,
		UpdateAt: user_mod.UpdateAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "Seuccess",
		"message": "User Registered",
		"data":    data,
	})
}
