package database

import (
	"fmt"
	"log"

	"github.com/faizulfikri/task-5-vix-btpn-Mohamad_Faizul_Fikri/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dbUser = "root"
const dbHost = "localhost"
const dbName = "data_users_btpn"
const dbPort = "3306"

var Db *gorm.DB

func initDB() *gorm.DB {
	Db = DatabaseConnection()
	return Db
}

// DatabaseConnection is creating new connection
func DatabaseConnection() *gorm.DB {
	var err error

	dsn := fmt.Sprintf("%s:WowLove@1234@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error Connecting to database mysql")
		log.Fatal(err)
	}

	db.Debug().AutoMigrate(&model.User{}, &model.Photo{})

	return db
}
