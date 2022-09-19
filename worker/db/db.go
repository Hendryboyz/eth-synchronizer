package db

import (
	"fmt"
	"log"

	"github.com/Hendryboyz/eth-synchronizer/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	err error
)

func Init() {
	connectionString := generateConnectionString()
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("DB connect fail")
	} else {
		log.Println("DB connected")
	}
}

func generateConnectionString() string {
	config := configs.GetConfig()
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetString("mysql.username"),
		config.GetString("mysql.pwd"),
		config.GetString("mysql.host"),
		config.GetInt32("mysql.port"),
		config.GetString("mysql.db"),
	)
}

func GetDBInstance() *gorm.DB {
	return db
}
