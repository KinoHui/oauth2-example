package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	MySQL struct {
		DataSource string
	}
	Redis redis.RedisConf
	Auth  struct {
		AccessSecret string
		AccessExpire int64
	}
	AutoApproveClients []string
}
