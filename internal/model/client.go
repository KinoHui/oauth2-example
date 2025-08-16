package model

import (
	"time"
)

// Client 客户端信息表
type Client struct {
	ID          string    `db:"id" json:"id"`                     // 客户端ID
	Secret      string    `db:"secret" json:"secret"`             // 客户端密钥
	Name        string    `db:"name" json:"name"`                 // 应用名称
	RedirectURL string    `db:"redirect_url" json:"redirect_url"` // 回调地址
	GrantType   string    `db:"grant_type" json:"grant_type"`     // 支持的授权模式
	Scope       string    `db:"scope" json:"scope"`               // 请求的权限范围
	CreatedAt   time.Time `db:"created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`     // 更新时间
}

// ClientRegisterReq 客户端注册请求
type ClientRegisterReq struct {
	Name        string `json:"name"`         // 应用名称
	RedirectURL string `json:"redirect_url"` // 回调地址
	GrantType   string `json:"grant_type"`   // 支持的授权模式
	Scope       string `json:"scope"`        // 请求的权限范围
}

// ClientRegisterResp 客户端注册响应
type ClientRegisterResp struct {
	ClientID     string `json:"client_id"`     // 客户端ID
	ClientSecret string `json:"client_secret"` // 客户端密钥
}
