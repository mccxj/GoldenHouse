package spider

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisManager struct {
	Client *redis.Client
	Prefix string
}

func (s *RedisManager) allUrlsKey() string {
	return s.Prefix + "_all_urls"
}

func (s *RedisManager) todoUrlsKey() string {
	return s.Prefix + "_todo_urls"
}

func (s *RedisManager) doingUrlsKey() string {
	return s.Prefix + "_doing_urls"
}

func (s *RedisManager) AppendUrl(url string) (bool, error) {
	c1 := s.Client.SIsMember(s.allUrlsKey(), url)
	if c1.Err() != nil {
		return false, c1.Err()
	}
	if c1.Val() {
		return false, nil
	}
	c2 := s.Client.LPush(s.todoUrlsKey(), url)
	if c2.Err() != nil {
		return false, c2.Err()
	}
	c3 := s.Client.SAdd(s.allUrlsKey(), url)
	if c3.Err() != nil {
		return false, c3.Err()
	}
	return true, nil
}

func (s *RedisManager) WaitUrl(timeout time.Duration) (string, error) {
	return s.Client.BRPopLPush(s.todoUrlsKey(), s.doingUrlsKey(), timeout).Result()
}

func (s *RedisManager) DoneUrl(url string) error {
	return s.Client.LRem(s.doingUrlsKey(), 0, url).Err()
}
