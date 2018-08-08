package spider

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Downloader interface {
}

func NewClient() *http.Client {
	httpProxy, ok := os.LookupEnv("http_proxy")
	var proxy func(*http.Request) (*url.URL, error)
	if ok {
		fmt.Println("using proxy: ", httpProxy)
		u := url.URL{}
		urlproxy, _ := u.Parse(httpProxy)
		proxy = http.ProxyURL(urlproxy)
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: proxy,
		},
	}
}

func NewRequest(url string) (*http.Request, error) {
	return http.NewRequest("GET", url, nil)
}
