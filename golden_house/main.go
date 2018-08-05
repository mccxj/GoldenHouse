package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	fmt.Println("start...")
	httpProxy, ok := os.LookupEnv("http_proxy")
	var proxy func(*http.Request) (*url.URL, error)
	if ok {
		urli := url.URL{}
		urlproxy, _ := urli.Parse(httpProxy)
		proxy = http.ProxyURL(urlproxy)
	}

	c := http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}

	url := "http://readfree.me/"
	req, err := http.NewRequest("GET", url, nil)
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
	req.Header.Add("Cookie", "sessionid=ffz8cqntoyjjqqjk9m3obdkxgliu23gz; csrftoken=VD0QbXiAlD2MjeE7ZNS9Zkgu9ivR0IdquUmqAX48rvPvn0M9UwZkJnGROEl3M9uA; Hm_lvt_375aa6d601368176e50751c1c6bf0e82=1533470014; Hm_lpvt_375aa6d601368176e50751c1c6bf0e82=1533475207")

	resp, err := c.Do(req)
	defer resp.Body.Close()
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exist := s.Attr("href"); exist {
			a := urlJoin(url, href)
			fmt.Println(a)
		}
	})
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
