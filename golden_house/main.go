package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mccxj/GoldenHouse/config"
	"github.com/mccxj/GoldenHouse/spider"
)

func main() {
	c := &config.SpiderConfig{}
	config.Load(c)
	fmt.Println(c)
	manager := &spider.RedisManager{
		Client: redis.NewClient(&redis.Options{
			Addr:     c.Redis.Addr,
			Password: c.Redis.Password,
			DB:       c.Redis.DB,
		}),
		Prefix: "readfree",
	}
	spider := spider.NewSpider(manager)
	spider.Run()
	fmt.Println("exit...")
}
