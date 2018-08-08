package spider

import (
	"net/url"
	"strings"
)

func UrlJoin(curr string, link string) string {
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

func UrlTrim(link string) string {
	li := strings.LastIndex(link, "#")
	if li == -1 {
		return link
	}
	return link[0:li]
}
