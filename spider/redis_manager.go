package spider

import (
	"github.com/go-redis/redis"
	"time"
)

const SITE_PREFIX = "readfree"
const ALL_URLS_KEY = SITE_PREFIX + "_all_urls"
const TODO_URLS_KEY = SITE_PREFIX + "_todo_urls"
const DOING_URLS_KEY = SITE_PREFIX + "_doing_urls"

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
