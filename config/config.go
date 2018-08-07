package config

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
)

type SpiderConfig struct {
	Redis redis
}

type redis struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

func Load(c *SpiderConfig) error {
	_, err := toml.DecodeFile("spider.toml", c)
	return err
}

func (sc SpiderConfig) String() string {
	body, _ := json.Marshal(sc)
	return string(body)
}
