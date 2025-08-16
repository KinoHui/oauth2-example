package types

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

// AuthorizeReq 授权请求
type AuthorizeReq struct {
	ClientID     string `form:"client_id"`     // 客户端ID
	ResponseType string `form:"response_type"` // 响应类型
	RedirectURI  string `form:"redirect_uri"`  // 重定向URI
	Scope        string `form:"scope"`         // 权限范围
	State        string `form:"state"`         // 状态参数
}

// AuthorizeResp 授权响应
type AuthorizeResp struct {
	Code  string `json:"code"`  // 授权码
	State string `json:"state"` // 状态参数
}

// TokenReq Token请求
type TokenReq struct {
	GrantType    string `form:"grant_type"`    // 授权类型
	Code         string `form:"code"`          // 授权码
	RedirectURI  string `form:"redirect_uri"`  // 重定向URI
	ClientID     string `form:"client_id"`     // 客户端ID
	ClientSecret string `form:"client_secret"` // 客户端密钥
}

// TokenResp Token响应
type TokenResp struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	TokenType    string `json:"token_type"`    // 令牌类型
	ExpiresIn    int64  `json:"expires_in"`    // 过期时间
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	Scope        string `json:"scope"`         // 权限范围
}

// UserInfoResp 用户信息响应
type UserInfoResp struct {
	UserID   string `json:"userid"`   // 用户ID
	Username string `json:"username"` // 用户名
	Phone    string `json:"phone"`    // 手机号
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// LoginResp 登录响应
type LoginResp struct {
	UserID string `json:"user_id"` // 用户ID
}

// ApproveReq 授权批准请求
type ApproveReq struct {
	ClientID string `json:"client_id"` // 客户端ID
	UserID   string `json:"user_id"`   // 用户ID
	Scope    string `json:"scope"`     // 权限范围
	Action   string `json:"action"`    // 动作：approve/reject
}
