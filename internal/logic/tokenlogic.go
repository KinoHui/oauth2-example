package logic

import (
	"context"
	"encoding/json"
	"errors"
	"oauth2-server/internal/svc"
	"oauth2-server/internal/types"
	"oauth2-server/internal/util"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type TokenLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenLogic {
	return &TokenLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TokenLogic) Token(req *types.TokenReq) (resp *types.TokenResp, err error) {
	// 验证授权类型
	if req.GrantType != "authorization_code" {
		return nil, errors.New("unsupported grant type")
	}

	// 验证客户端
	client, err := l.svcCtx.ClientModel.FindByID(l.ctx, req.ClientID)
	if err != nil {
		return nil, errors.New("invalid client")
	}

	// 验证客户端密钥
	if client.Secret != req.ClientSecret {
		return nil, errors.New("invalid client secret")
	}

	// 从Redis获取授权码数据
	redisStore := util.NewRedisStore(l.svcCtx.Redis)
	codeDataStr, err := redisStore.GetCode(l.ctx, req.Code)
	if err != nil {
		return nil, errors.New("invalid authorization code")
	}

	// 解析授权码数据
	var codeData map[string]interface{}
	err = json.Unmarshal([]byte(codeDataStr), &codeData)
	if err != nil {
		return nil, errors.New("invalid authorization code data")
	}

	// 验证客户端ID
	if codeData["client_id"] != req.ClientID {
		return nil, errors.New("client_id mismatch")
	}

	// 验证重定向URI
	if codeData["redirect_uri"] != req.RedirectURI {
		return nil, errors.New("redirect_uri mismatch")
	}

	// 生成访问令牌
	accessToken, err := util.GenerateToken(
		codeData["user_id"].(string),
		req.ClientID,
		codeData["scope"].(string),
		l.svcCtx.Config.Auth.AccessSecret,
		l.svcCtx.Config.Auth.AccessExpire,
	)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken := uuid.New().String()

	// 存储访问令牌到Redis
	tokenData := map[string]interface{}{
		"user_id":   codeData["user_id"],
		"client_id": req.ClientID,
		"scope":     codeData["scope"],
	}
	err = redisStore.StoreAccessToken(l.ctx, accessToken, tokenData, time.Duration(l.svcCtx.Config.Auth.AccessExpire)*time.Second)
	if err != nil {
		return nil, err
	}

	// 存储刷新令牌到Redis
	refreshTokenData := map[string]interface{}{
		"user_id":   codeData["user_id"],
		"client_id": req.ClientID,
		"scope":     codeData["scope"],
	}
	err = redisStore.StoreRefreshToken(l.ctx, refreshToken, refreshTokenData, 30*24*time.Hour) // 30天
	if err != nil {
		return nil, err
	}

	// 删除授权码
	redisStore.DeleteCode(l.ctx, req.Code)

	return &types.TokenResp{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    l.svcCtx.Config.Auth.AccessExpire,
		RefreshToken: refreshToken,
		Scope:        codeData["scope"].(string),
	}, nil
}
