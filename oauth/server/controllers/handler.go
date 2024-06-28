package controllers

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type controller struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func New(db *gorm.DB, rdb *redis.Client) controller {
	return controller{db, rdb}
}
