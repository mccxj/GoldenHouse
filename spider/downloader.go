package spider

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

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
