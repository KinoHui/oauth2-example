package logic

import (
	"context"
	"oauth2-server/internal/model"
	"oauth2-server/internal/svc"
	"oauth2-server/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClientRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClientRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientRegisterLogic {
	return &ClientRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClientRegisterLogic) ClientRegister(req *types.ClientRegisterReq) (resp *types.ClientRegisterResp, err error) {
	// 创建客户端记录
	client := &model.Client{
		Name:        req.Name,
		RedirectURL: req.RedirectURL,
		GrantType:   req.GrantType,
		Scope:       req.Scope,
	}

	// 插入数据库
	_, err = l.svcCtx.ClientModel.Insert(l.ctx, client)
	if err != nil {
		return nil, err
	}

	return &types.ClientRegisterResp{
		ClientID:     client.ID,
		ClientSecret: client.Secret,
	}, nil
}
