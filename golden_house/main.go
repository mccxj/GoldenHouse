package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const ENTRANCE string = "http://readfree.me/"

var regs = make([]*regexp.Regexp, 0, 10)

func init() {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(`^http://readfree.me/\?page=\d+$`)
	regs = append(regs, reg)
	reg, _ = regexp.Compile(`^http://readfree.me/book/\d+/$`)
	regs = append(regs, reg)
}

func main() {
	urls := make(map[string]bool, 100)
	urls[ENTRANCE] = false

	queue := make([]string, 0)
	// Push to the queue
	queue = append(queue, ENTRANCE)

	httpProxy, ok := os.LookupEnv("http_proxy")
	var proxy func(*http.Request) (*url.URL, error)
	if ok {
		fmt.Println("using proxy: ", httpProxy)
		urli := url.URL{}
		urlproxy, _ := urli.Parse(httpProxy)
		proxy = http.ProxyURL(urlproxy)
	}

	c := http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}

	fmt.Println("start...")
	for len(queue) != 0 {
		time.Sleep(time.Second * 1)
		todoUrl := queue[0]
		fmt.Println(todoUrl)

		req, err := http.NewRequest("GET", todoUrl, nil)
		if err != nil {
			log.Fatalln(err)
		}

		req.Header.Add("Host", "readfree.me")
		req.Header.Add("Connection", "keep-alive")
		req.Header.Add("Cache-Control", "max-age=0")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.84 Safari/537.36")
		req.Header.Add("Upgrade-Insecure-Requests", "1")
		req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		//req.Header.Add("Accept-Encoding", "gzip, deflate")
		req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
		req.Header.Add("Cookie", "Hm_lvt_375aa6d601368176e50751c1c6bf0e82=1532316147,1532316183,1533547490; sessionid=8ukyqbsdz2bnfrarime2tzcg4u12e13t; csrftoken=OQMt4Lunj58NS0bVGyrmdzopuM70lByBYsHdBLIG2eKAkg0tLTJjXd3coac9g5AP; Hm_lpvt_375aa6d601368176e50751c1c6bf0e82=1533548656")

		resp, err := c.Do(req)
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
		}

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			if href, exist := s.Attr("href"); exist {
				newUrl := urlTrim(urlJoin(todoUrl, href))
				if strings.HasPrefix(newUrl, ENTRANCE) && isValidUrl(newUrl) {
					if _, exist := urls[newUrl]; !exist {
						urls[newUrl] = false
						queue = append(queue, newUrl)
					}
				}
			}
		})

		queue = queue[1:]
	}
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
