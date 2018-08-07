package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mccxj/GoldenHouse/spider"
)

func main() {
	manager := &spider.RedisManager{redis.NewClient(&redis.Options{
		Addr:     "100.101.120.200:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
	spider := spider.NewSpider(manager)
	spider.Run()
	fmt.Println("exit...")
}
