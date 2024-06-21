package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bigbluewhale111/rest_api/models"
)

func Init() *gorm.DB {
	dbURL := "postgres://pg:pg_pass123@localhost:5433/task"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Task{})

	return db
}
