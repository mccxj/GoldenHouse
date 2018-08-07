package spider

import (
	"net/url"
	"regexp"
	"strings"
)

var regs = make([]*regexp.Regexp, 0, 10)

func init() {
	var reg *regexp.Regexp
	reg, _ = regexp.Compile(`^http://readfree.me/\?page=\d+$`)
	regs = append(regs, reg)
	reg, _ = regexp.Compile(`^http://readfree.me/book/\d+/$`)
	regs = append(regs, reg)
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
