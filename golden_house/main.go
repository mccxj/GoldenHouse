package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
	"github.com/pingcap/tidb/util/goroutine_pool"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const ENTRANCE string = "http://readfree.me/"
const SITE_PREFIX = "readfree"
const ALL_URLS_KEY = SITE_PREFIX + "_all_urls"
const TODO_URLS_KEY = SITE_PREFIX + "_todo_urls"
const DOING_URLS_KEY = SITE_PREFIX + "_doing_urls"

var regs = make([]*regexp.Regexp, 0, 10)

type Manager interface {
	AppendUrl(url string) (bool, error)
	WaitUrl(timeout time.Duration) (string, error)
	DoneUrl(string) error
}

type RedisManager struct {
	client *redis.Client
}

func (s *RedisManager) AppendUrl(url string) (bool, error) {
	c1 := s.client.SIsMember(ALL_URLS_KEY, url)
	if c1.Err() != nil {
		return false, c1.Err()
	}
	if c1.Val() {
		return false, nil
	}
	c2 := s.client.LPush(TODO_URLS_KEY, url)
	if c2.Err() != nil {
		return false, c2.Err()
	}
	c3 := s.client.SAdd(ALL_URLS_KEY, url)
	if c3.Err() != nil {
		return false, c3.Err()
	}
	return true, nil
}

func (s *RedisManager) WaitUrl(timeout time.Duration) (string, error) {
	return s.client.BRPopLPush(TODO_URLS_KEY, DOING_URLS_KEY, timeout).Result()
}

func (s *RedisManager) DoneUrl(url string) error {
	return s.client.LRem(DOING_URLS_KEY, 0, url).Err()
}

type Spider struct {
	Manager
}

func NewSpider(manager Manager) *Spider {
	return &Spider{
		Manager: manager,
	}
}

func (s *Spider) Run() {
	pool := gp.New(200 * time.Millisecond)
	for i := 0; i < 10; i++ {
		pool.Go(func() {
			c := NewClient()
			for {
				time.Sleep(time.Second * 2)
				todoUrl, _ := s.WaitUrl(time.Second * 10)
				fmt.Println("=>", todoUrl)
				req, err := NewRequest(todoUrl)
				if err != nil {
					log.Fatalln(err)
				}
				resp, err := c.Do(req)
				if err != nil {
					log.Fatalln(err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
				}

				doc, err := goquery.NewDocumentFromReader(resp.Body)
				if err != nil {
					log.Fatal(err)
				}

				doc.Find("a").Each(func(i int, gs *goquery.Selection) {
					if href, exist := gs.Attr("href"); exist {
						newUrl := urlTrim(urlJoin(todoUrl, href))
						if strings.HasPrefix(newUrl, ENTRANCE) && isValidUrl(newUrl) {
							s.AppendUrl(newUrl)
						}
					}
				})
				s.DoneUrl(todoUrl)
			}
		})
	}

	s.AppendUrl(ENTRANCE)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func init() {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(`^http://readfree.me/\?page=\d+$`)
	regs = append(regs, reg)
	reg, _ = regexp.Compile(`^http://readfree.me/book/\d+/$`)
	regs = append(regs, reg)
}

func NewClient() *http.Client {
	httpProxy, ok := os.LookupEnv("http_proxy")
	var proxy func(*http.Request) (*url.URL, error)
	if ok {
		fmt.Println("using proxy: ", httpProxy)
		urli := url.URL{}
		urlproxy, _ := urli.Parse(httpProxy)
		proxy = http.ProxyURL(urlproxy)
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}
}

func NewRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return req, err
	}

	req.Header.Add("Host", "readfree.me")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Add("Cookie", "Hm_lvt_375aa6d601368176e50751c1c6bf0e82=1532316147,1532316183,1533547490; sessionid=8ukyqbsdz2bnfrarime2tzcg4u12e13t; csrftoken=OQMt4Lunj58NS0bVGyrmdzopuM70lByBYsHdBLIG2eKAkg0tLTJjXd3coac9g5AP; Hm_lpvt_375aa6d601368176e50751c1c6bf0e82=1533621082")
	return req, nil
}

func main() {
	manager := &RedisManager{redis.NewClient(&redis.Options{
		Addr:     "100.101.120.200:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})}
	spider := NewSpider(manager)
	spider.Run()
	fmt.Println("exit...")
}

func urlJoin(curr string, link string) string {
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		return link
	}

	u, _ := url.Parse(curr)
	host := u.Scheme + "://" + u.Host
	if strings.HasPrefix(link, "/") {
		return host + link
	}
	li := strings.LastIndex(curr, "/")
	if li == len(u.Scheme+"://")-1 {
		return host + "/" + link
	}
	return curr[0:li] + "/" + link
}

func urlTrim(link string) string {
	li := strings.LastIndex(link, "#")
	if li == -1 {
		return link
	}
	return link[0:li]
}

func isValidUrl(link string) bool {
	for _, reg := range regs {
		if reg.MatchString(link) {
			return true
		}
	}
	return false
}
