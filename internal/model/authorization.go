package model

import (
	"time"
)

// Authorization 权限申请记录表
type Authorization struct {
	ID        int64     `db:"id" json:"id"`                 // 主键ID
	ClientID  string    `db:"client_id" json:"client_id"`   // 客户端ID
	UserID    string    `db:"user_id" json:"user_id"`       // 用户ID
	Scope     string    `db:"scope" json:"scope"`           // 请求的权限范围
	Code      string    `db:"code" json:"code"`             // 授权码
	Status    string    `db:"status" json:"status"`         // 状态：pending/approved/rejected
	CreatedAt time.Time `db:"created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // 更新时间
}

// AuthorizationReq 授权请求
type AuthorizationReq struct {
	ClientID     string `json:"client_id"`     // 客户端ID
	ResponseType string `json:"response_type"` // 响应类型
	RedirectURI  string `json:"redirect_uri"`  // 重定向URI
	Scope        string `json:"scope"`         // 权限范围
	State        string `json:"state"`         // 状态参数
}

// AuthorizationResp 授权响应
type AuthorizationResp struct {
	Code  string `json:"code"`  // 授权码
	State string `json:"state"` // 状态参数
}

// TokenReq Token请求
type TokenReq struct {
	GrantType    string `json:"grant_type"`    // 授权类型
	Code         string `json:"code"`          // 授权码
	RedirectURI  string `json:"redirect_uri"`  // 重定向URI
	ClientID     string `json:"client_id"`     // 客户端ID
	ClientSecret string `json:"client_secret"` // 客户端密钥
}

// TokenResp Token响应
type TokenResp struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	TokenType    string `json:"token_type"`    // 令牌类型
	ExpiresIn    int64  `json:"expires_in"`    // 过期时间
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	Scope        string `json:"scope"`         // 权限范围
}

// UserInfo 用户信息
type UserInfo struct {
	UserID   string `json:"userid"`   // 用户ID
	Username string `json:"username"` // 用户名
	Phone    string `json:"phone"`    // 手机号
}
