package logic

import (
	"context"
	"errors"
	"log"
	"oauth2-server/internal/model"
	"oauth2-server/internal/svc"
	"oauth2-server/internal/types"
	"oauth2-server/internal/util"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthorizeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthorizeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthorizeLogic {
	return &AuthorizeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthorizeLogic) Authorize(req *types.AuthorizeReq) (resp *types.AuthorizeResp, err error) {
	// 验证客户端
	client, err := l.svcCtx.ClientModel.FindByID(l.ctx, req.ClientID)
	if err != nil {
		return nil, errors.New("invalid client")
	}
	log.Println(client)

	// 验证重定向URI
	if client.RedirectURL != req.RedirectURI {
		return nil, errors.New("invalid redirect uri")
	}

	// 验证响应类型
	if req.ResponseType != "code" {
		return nil, errors.New("unsupported response type")
	}

	// 检查是否为自动批准的客户端
	isAutoApprove := false
	for _, autoClientID := range l.svcCtx.Config.AutoApproveClients {
		if autoClientID == req.ClientID {
			isAutoApprove = true
			break
		}
	}

	// 如果是自动批准的客户端，直接生成授权码
	if isAutoApprove {
		return l.generateAuthorizationCode(req, client, "test_user") // 使用默认用户
	}

	// 否则需要用户登录和授权
	// 这里应该重定向到登录页面，暂时返回错误
	return nil, errors.New("user authorization required")
}

func (l *AuthorizeLogic) generateAuthorizationCode(req *types.AuthorizeReq, client *model.Client, userID string) (*types.AuthorizeResp, error) {
	// 创建授权记录
	auth := &model.Authorization{
		ClientID: req.ClientID,
		UserID:   userID,
		Scope:    req.Scope,
		Status:   "approved",
	}

	// 插入数据库
	_, err := l.svcCtx.AuthorizationModel.Insert(l.ctx, auth)
	if err != nil {
		return nil, err
	}

	// 存储授权码到Redis
	redisStore := util.NewRedisStore(l.svcCtx.Redis)
	codeData := map[string]interface{}{
		"client_id":    req.ClientID,
		"user_id":      userID,
		"scope":        req.Scope,
		"redirect_uri": req.RedirectURI,
	}

	err = redisStore.StoreCode(l.ctx, auth.Code, codeData, 10*time.Minute)
	if err != nil {
		return nil, err
	}

	return &types.AuthorizeResp{
		Code:  auth.Code,
		State: req.State,
	}, nil
}
