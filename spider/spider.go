package spider

import (
	"fmt"
	"github.com/pingcap/tidb/util/goroutine_pool"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Site interface {
	Entrance() string
	CustomReq(req *http.Request)
	Extract(todoUrl string, body io.ReadCloser) (urls []string)
}

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

func (s *Spider) Run(site Site) {
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
				site.CustomReq(req)
				resp, err := c.Do(req)
				if err != nil {
					log.Fatalln(err)
				}
				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
				}

				newUrls := site.Extract(todoUrl, resp.Body)
				for _, newUrl := range newUrls {
					s.AppendUrl(newUrl)
				}
				s.DoneUrl(todoUrl)
			}
		})
	}

	_, err := s.AppendUrl(site.Entrance())
	if err != nil {
		log.Fatalln(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
