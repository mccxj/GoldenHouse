package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/pingcap/tidb/util/goroutine_pool"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const ENTRANCE string = "http://readfree.me/"

type Manager interface {
	AppendUrl(url string) (bool, error)
	WaitUrl(timeout time.Duration) (string, error)
	DoneUrl(string) error
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
