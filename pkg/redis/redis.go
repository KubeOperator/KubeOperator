package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

var Client *redis.Client

const phaseName = "redis"

type InitRedisPhase struct {
	Host       string
	Port       int
	DB         int
	MaxRetries int
}

func (i *InitRedisPhase) Init() error {
	Client = redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", i.Host, i.Port),
		DB:         i.DB,
		MaxRetries: i.MaxRetries,
	})
	_, err := Client.Ping().Result()
	return err
}

func (i *InitRedisPhase) PhaseName() string {
	return phaseName
}
