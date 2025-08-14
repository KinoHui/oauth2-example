package main

import (
	"time"
)

// ClientInfo 客户端信息表
type ClientInfo struct {
	ID          string    `json:"id" db:"id"`
	Secret      string    `json:"secret" db:"secret"`
	Name        string    `json:"name" db:"name"`
	RedirectURL string    `json:"redirect_url" db:"redirect_url"`
	GrantType   string    `json:"grant_type" db:"grant_type"`     // 支持的授权模式
	Scope       string    `json:"scope" db:"scope"`               // 请求的权限范围
	AutoApprove bool      `json:"auto_approve" db:"auto_approve"` // 是否自动授权
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// AuthRecord 权限申请记录表
type AuthRecord struct {
	ID        int64     `json:"id" db:"id"`
	ClientID  string    `json:"client_id" db:"client_id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Scope     string    `json:"scope" db:"scope"`
	Approved  bool      `json:"approved" db:"approved"` // 用户是否同意授权
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ClientRegistrationRequest 客户端注册请求
type ClientRegistrationRequest struct {
	Name        string `json:"name" binding:"required"`
	RedirectURL string `json:"redirect_url" binding:"required"`
	GrantType   string `json:"grant_type" binding:"required"`
	Scope       string `json:"scope" binding:"required"`
}

// ClientRegistrationResponse 客户端注册响应
type ClientRegistrationResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Message      string `json:"message"`
}

// UserInfo 用户信息
type UserInfo struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}
