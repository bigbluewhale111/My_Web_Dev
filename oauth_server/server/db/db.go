package db

import (
	"log"
	"os"

	"github.com/bigbluewhale111/oauth_server/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("OAUTH_DB")), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.User{}, &models.AccessToken{})

	return db
}
