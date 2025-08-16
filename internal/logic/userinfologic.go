package logic

import (
	"context"
	"errors"
	"net/http"
	"oauth2-server/internal/svc"
	"oauth2-server/internal/types"
	"oauth2-server/internal/util"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(r *http.Request) (resp *types.UserInfoResp, err error) {
	// 从请求头获取Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	// 解析Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid authorization header")
	}

	token := parts[1]

	// 验证JWT token
	claims, err := util.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// 从Redis验证token是否有效
	redisStore := util.NewRedisStore(l.svcCtx.Redis)
	_, err = redisStore.GetAccessToken(l.ctx, token)
	if err != nil {
		return nil, errors.New("token expired or invalid")
	}

	// 根据scope返回相应的用户信息
	userInfo := &types.UserInfoResp{}

	// 模拟用户数据，实际应该从数据库获取
	if strings.Contains(claims.Scope, "userid") {
		userInfo.UserID = claims.UserID
	}

	if strings.Contains(claims.Scope, "profile") {
		// 模拟用户信息
		userInfo.Username = "test_user"
		userInfo.Phone = "13800138000"
	}

	return userInfo, nil
}
