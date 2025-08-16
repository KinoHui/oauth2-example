package util

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

// RedisStore Redis存储工具
type RedisStore struct {
	redis redis.Redis
}

// NewRedisStore 创建Redis存储实例
func NewRedisStore(r redis.Redis) *RedisStore {
	return &RedisStore{redis: r}
}

// StoreCode 存储授权码
func (rs *RedisStore) StoreCode(ctx context.Context, code string, data interface{}, expire time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := "oauth:code:" + code
	return rs.redis.SetexCtx(ctx, key, string(jsonData), int(expire.Seconds()))
}

// GetCode 获取授权码数据
func (rs *RedisStore) GetCode(ctx context.Context, code string) (string, error) {
	key := "oauth:code:" + code
	return rs.redis.GetCtx(ctx, key)
}

// DeleteCode 删除授权码
func (rs *RedisStore) DeleteCode(ctx context.Context, code string) error {
	key := "oauth:code:" + code
	_, err := rs.redis.DelCtx(ctx, key)
	return err
}

// StoreAccessToken 存储访问令牌
func (rs *RedisStore) StoreAccessToken(ctx context.Context, token string, data interface{}, expire time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := "oauth:token:" + token
	return rs.redis.SetexCtx(ctx, key, string(jsonData), int(expire.Seconds()))
}

// GetAccessToken 获取访问令牌数据
func (rs *RedisStore) GetAccessToken(ctx context.Context, token string) (string, error) {
	key := "oauth:token:" + token
	return rs.redis.GetCtx(ctx, key)
}

// DeleteAccessToken 删除访问令牌
func (rs *RedisStore) DeleteAccessToken(ctx context.Context, token string) error {
	key := "oauth:token:" + token
	_, err := rs.redis.DelCtx(ctx, key)
	return err
}

// StoreRefreshToken 存储刷新令牌
func (rs *RedisStore) StoreRefreshToken(ctx context.Context, refreshToken string, data interface{}, expire time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := "oauth:refresh:" + refreshToken
	return rs.redis.SetexCtx(ctx, key, string(jsonData), int(expire.Seconds()))
}

// GetRefreshToken 获取刷新令牌数据
func (rs *RedisStore) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := "oauth:refresh:" + refreshToken
	return rs.redis.GetCtx(ctx, key)
}

// DeleteRefreshToken 删除刷新令牌
func (rs *RedisStore) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := "oauth:refresh:" + refreshToken
	_, err := rs.redis.DelCtx(ctx, key)
	return err
}
