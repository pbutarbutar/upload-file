package config

import (
	"brks/domain"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var db *gorm.DB
var port string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(err).Msg("unable to load env through env config")
	}
}

func SetupModels() {
	dsn := os.Getenv("DSN_DB")
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return
	}

	// create table if it does not exist
	if !db.HasTable(&domain.Upload{}) {
		db.CreateTable(&domain.Upload{})
	}

	if err != nil {
		fmt.Println("Failed to connect to database!", err)
		return
	}

	SetUpDBConnection(db)
	//Set Port
	SetPortConnection(os.Getenv("PORT"))
}

func SetUpDBConnection(DB *gorm.DB) {
	db = DB
}

func GetDBConnection() *gorm.DB {
	return db
}

func SetPortConnection(Port string) {
	port = Port
}

func GetPortConnection() string {
	return port
}
