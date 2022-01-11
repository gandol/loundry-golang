package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"loundry/api/src/helper"
	"loundry/api/src/models"
	"os"
)

var (
	DBConn *gorm.DB
)

func ConnectDb() {
	usernameDB := helper.ReadEnv("DB_USERNAME")
	passDB := helper.ReadEnv("DB_PASSWORD")
	hostDB := helper.ReadEnv("DB_HOST")
	portDB := helper.ReadEnv("DB_PORT")
	dbName := helper.ReadEnv("DB_NAME")
	connection, err := gorm.Open(mysql.Open(usernameDB+":"+passDB+"@tcp("+hostDB+":"+portDB+")/"+dbName+"?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		log.Fatal("unable to connect to database")
		os.Exit(2)

	}
	connection.AutoMigrate(&models.Users{})
	connection.AutoMigrate(&models.Customer{})
	connection.AutoMigrate(&models.Transaksi{})
	connection.AutoMigrate(&models.TransaksiItems{})
	connection.AutoMigrate(&models.Notifications{})
	connection.AutoMigrate(&models.Items{})
	DBConn = connection
}
