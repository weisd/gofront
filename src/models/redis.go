package models

import (
	"github.com/garyburd/redigo/redis"
	"modules/log"
	"time"
)

type RedisService struct {
	Enable bool
	Addr   string
	Passwd string

	MaxIdle     int
	IdleTimeout int
}

var RedisPools map[string]*redis.Pool

func Redis(name ...string) *redis.Pool {
	k := "default"
	if len(name) > 0 {
		k = name[0]
	}

	if pool, ok := RedisPools[k]; ok {
		return pool

	}

	log.Fatal("unkown redis %s", k)

	return nil
}

func InitRedis(confs map[string]RedisService) {
	RedisPools = make(map[string]*redis.Pool)
	for name, conf := range confs {
		RedisPools[name] = newRedis(conf)
	}

	log.Debug("init redis done %v", RedisPools)
}

func newRedis(conf RedisService) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: time.Duration(conf.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", conf.Addr)
			if err != nil {
				return nil, err
			}
			if len(conf.Passwd) > 0 {
				if _, err := c.Do("AUTH", conf.Passwd); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
