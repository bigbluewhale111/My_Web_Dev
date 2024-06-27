package cache

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func Init() *redis.Client {
	opt, err := redis.ParseURL(os.Getenv("REDIS"))
	if err != nil {
		log.Fatalln(err)
	}
	rdb := redis.NewClient(opt)

	return rdb
}
