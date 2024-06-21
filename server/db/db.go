package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bigbluewhale111/rest_api/models"
)

func Init() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB")), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Task{})

	return db
}
