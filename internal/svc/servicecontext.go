package svc

import (
	"oauth2-server/internal/config"
	"oauth2-server/internal/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config             config.Config
	DB                 sqlx.SqlConn
	Redis              redis.Redis
	ClientModel        model.ClientModel
	AuthorizationModel model.AuthorizationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MySQL.DataSource)

	return &ServiceContext{
		Config:             c,
		DB:                 conn,
		Redis:              *redis.MustNewRedis(c.Redis),
		ClientModel:        model.NewClientModel(conn),
		AuthorizationModel: model.NewAuthorizationModel(conn),
	}
}
