package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
	"github.com/mccxj/GoldenHouse/config"
	"github.com/mccxj/GoldenHouse/spider"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var regs = make([]*regexp.Regexp, 0, 10)

func init() {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(`^http://readfree.me/\?page=\d$`)
	regs = append(regs, reg)
	reg, _ = regexp.Compile(`^http://readfree.me/book/\d/$`)
	regs = append(regs, reg)
}

func isValidUrl(link string) bool {
	for _, reg := range regs {
		if reg.MatchString(link) {
			return true
		}
	}
	return false
}

type ReadfreeSite struct {
	Headers map[string]string
}

func (site *ReadfreeSite) Entrance() string {
	return "http://readfree.me/"
}

func (site *ReadfreeSite) CustomReq(req *http.Request) {
	for k, v := range site.Headers {
		req.Header.Add(k, v)
	}
}

func (site *ReadfreeSite) Extract(todoUrl string, body io.ReadCloser) (urls []string) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	newUrls := make([]string, 0, 10)
	doc.Find("a").Each(func(i int, gs *goquery.Selection) {
		if href, exist := gs.Attr("href"); exist {
			newUrl := spider.UrlTrim(spider.UrlJoin(todoUrl, href))
			if strings.HasPrefix(newUrl, site.Entrance()) && isValidUrl(newUrl) {
				newUrls = append(newUrls, newUrl)
			}
		}
	})
	return newUrls
}

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
	spider.Run(&ReadfreeSite{
		Headers: c.Headers,
	})
	fmt.Println("exit...")
}
