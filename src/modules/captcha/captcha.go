package captcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dchest/captcha"
	"github.com/garyburd/redigo/redis"
	"github.com/labstack/echo"
	"models"
	"modules/log"
	"path"
	"time"

	"strings"
)

func InitCaptcha(op ...Options) error {
	option := Options{}
	if len(op) > 0 {
		option = op[0]
	}

	option = PreOption(option)

	switch option.Driver {
	case "memory":
		return nil
	case "redis":

		var r RedistoreOption

		if err := json.Unmarshal([]byte(option.Config), &r); err != nil {
			return err
		}

		r.Expire = option.Exprie

		store, err := NewRedisStore(r)
		if err != nil {
			return err
		}
		captcha.SetCustomStore(store)
		return nil
	}

	return fmt.Errorf("Unkown Driver %s", option.Driver)

}

type Options struct {
	Driver string
	Config string
	Exprie int
}

func PreOption(op ...Options) Options {
	option := Options{}
	if len(op) > 0 {
		option = op[0]
	}

	defaultOption := Options{"memory", "", 30}

	if len(option.Driver) > 0 {
		defaultOption.Driver = strings.ToLower(option.Driver)
	}

	if len(option.Config) > 0 {
		defaultOption.Config = option.Driver
	}

	if option.Exprie > 0 {
		defaultOption.Exprie = option.Exprie
	}

	return defaultOption

}

func Server() echo.HandlerFunc {
	return func(c echo.Context) error {

		fileName := c.P(0)

		ext := path.Ext(fileName)
		id := fileName[:len(fileName)-len(ext)]

		var content bytes.Buffer

		if c.FormValue("reload") != "" {
			captcha.Reload(id)
		}

		captcha.WriteImage(&content, id, captcha.StdWidth, captcha.StdHeight)

		c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "0")
		c.Response().Header().Set("Content-Type", "image/png")

		return c.ServeContent(bytes.NewReader(content.Bytes()), id+ext, time.Now())

	}
}

/////////////// redis store /////////////////

type RedisStore struct {
	Prefix string
	Redis  *redis.Pool
	Expire int
}

func (s *RedisStore) Set(id string, digits []byte) {
	conn := s.Redis.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", fmt.Sprintf("%s.%s", id), s.Expire, digits)
	if err != nil {
		log.Error("captcha RedisStore set err : %v", err)
	}
}

func (s *RedisStore) Get(id string, clear bool) (digits []byte) {
	conn := s.Redis.Get()
	defer conn.Close()
	var err error
	digits, err = redis.Bytes(conn.Do("GET", fmt.Sprintf("%s.%s", id)))
	if err != nil {
		log.Error("captcha RedisStore set err : %v", err)
		return
	}

	if clear {
		conn.Do("DEL", fmt.Sprintf("%s.%s", id))
	}

	return
}

type RedistoreOption struct {
	RedisPrefix string
	RedisName   string
	Expire      int
}

func NewRedisStore(op RedistoreOption) (captcha.Store, error) {
	if !models.HasRedis(op.RedisName) {
		return nil, fmt.Errorf("Unkown redis %s", op.RedisName)
	}
	// @todo config able
	return &RedisStore{Prefix: op.RedisPrefix, Redis: models.Redis(op.RedisName), Expire: op.Expire}, nil
}
