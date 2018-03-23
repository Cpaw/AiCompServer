package db

import (
	// _ "github.com/go-sql-driver/mysql"
	"AiCompServer/app/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/mattn/go-sqlite3"
	"github.com/revel/revel"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitDB() {
	dbhost := os.Getenv("DBHOST")
	// dbname := os.Getenv("DBNAME")
	// dbuser := os.Getenv("DBUSER")
	// dbpass := os.Getenv("DBPASSWORD")
	// db, err := gorm.Open("sqlite3", dbInfoString())
	db, err := gorm.Open("postgres", "host="+dbhost+" port=5432 user=gorm dbname=gorm sslmode=disable password=yatuhashi-api")
	for {
		if err != nil {
			log.Println("Failed to connect to database: %v\n", err)
		} else {
			break
		}
		time.Sleep(1 * time.Second)
		db, err := gorm.Open("postgres", "host="+dbhost+" port=5432 user=gorm dbname=gorm sslmode=disable password=yatuhashi-api")
	}

	db.DB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Challenge{})
	db.AutoMigrate(&models.Answer{})
	DB = db
	// DB.Create(&models.User{Username: "tester"})
}

func dbInfoString() string {
	s, b := revel.Config.String("db.info")
	if !b {
		log.Panicf("database info not found")
	}
	return s
}
