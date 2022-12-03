package main

import (
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/database"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/model"
	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/router"
)

func main() {
	db := database.DatabaseConnection()
	db.AutoMigrate(&model.User{})

	r := router.SetUpRouter(db)
	r.Run(":8080")
}
